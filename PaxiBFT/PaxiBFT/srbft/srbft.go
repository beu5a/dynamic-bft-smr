package srbft

import (
	"crypto/md5"
	"time"

	"github.com/salemmohammed/PaxiBFT"
	"github.com/salemmohammed/PaxiBFT/log"
)

type status int8

const (
	NONE status = iota
	PREPREPARED
	PREPARED
	COMMITTED
	RECEIVED
)

// add special request ????

// log's entries
type entry struct {
	ballot    PaxiBFT.Ballot
	view      PaxiBFT.View
	command   PaxiBFT.Command
	commit    bool
	active    bool
	Leader    bool
	request   *PaxiBFT.Request
	timestamp time.Time
	Digest    []byte
	Q1        *PaxiBFT.Quorum
	Q2        *PaxiBFT.Quorum
	Q3        *PaxiBFT.Quorum
	Q4        *PaxiBFT.Quorum
	Pstatus   status
	Cstatus   status
	Rstatus   status
}

// srbft instance
type srbft struct {
	PaxiBFT.Node

	config []PaxiBFT.ID
	N      PaxiBFT.Config
	log    map[int]*entry // log ordered by slot

	slot               int            // highest slot number
	view               PaxiBFT.View   // view number
	ballot             PaxiBFT.Ballot // highest ballot number
	execute            int            // next execute slot number
	requests           []*PaxiBFT.Request
	quorum             *PaxiBFT.Quorum // phase 1 quorum
	ReplyWhenCommit    bool
	RecivedReq         bool
	Member             *PaxiBFT.Memberlist
	strugglerThreshold int
}

// Newsrbft creates new srbft instance
func Newsrbft(n PaxiBFT.Node, options ...func(*srbft)) *srbft {
	p := &srbft{
		Node: n,
		log:  make(map[int]*entry, PaxiBFT.GetConfig().BufferSize),

		quorum: PaxiBFT.NewQuorum(),
		slot:   -1,

		requests:           make([]*PaxiBFT.Request, 0),
		ReplyWhenCommit:    false,
		RecivedReq:         false,
		Member:             PaxiBFT.NewMember(),
		strugglerThreshold: 3,
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// Digest message
func GetMD5Hash(r *PaxiBFT.Request) []byte {
	hasher := md5.New()
	hasher.Write([]byte(r.Command.Value))
	return []byte(hasher.Sum(nil))
}

func (p *srbft) HandleRequest(r PaxiBFT.Request, s int) {
	log.Debugf("<--------------------HandleRequest------------------>")

	e := p.log[s]
	e.Digest = GetMD5Hash(&r)
	log.Debugf("[p.ballot.ID %v, p.ballot %v ]", p.ballot.ID(), p.ballot)
	log.Debugf("PrePrepare will be called")
	p.PrePrepare(&r, &e.Digest, s)
	p.LightPrePrepare(&r, s)
}

// Pre_prepare starts phase 1 PrePrepare
// the primary will send <<pre-prepare,v,n,d(m)>,m>
func (p *srbft) PrePrepare(r *PaxiBFT.Request, s *[]byte, slt int) {
	log.Debugf("<--------------------PrePrepare------------------>")

	_, ok := p.log[p.slot]
	if !ok {
		p.log[p.slot] = &entry{
			ballot:    p.ballot,
			view:      p.view,
			command:   r.Command,
			commit:    false,
			active:    false,
			Leader:    false,
			request:   r,
			timestamp: time.Now(),
			Digest:    GetMD5Hash(r),
			Q1:        PaxiBFT.NewQuorum(),
			Q2:        PaxiBFT.NewQuorum(),
			Q3:        PaxiBFT.NewQuorum(),
			Q4:        PaxiBFT.NewQuorum(),
		}
	}

	for id := 1; id < p.strugglerThreshold; id++ {
		if id != p.ID().Node() {

			p.Send(PaxiBFT.NewID(1, id), PrePrepare{
				Ballot:  p.ballot,
				ID:      p.ID(),
				View:    p.view,
				Slot:    slt,
				Request: *r,
				Digest:  *s,
				Command: r.Command,
			})
			log.Debugf("PrePrepare sent to %v", id)
			log.Debugf("Struggler Threshold = %v", p.strugglerThreshold)
		}
	}
	log.Debugf("++++++ PrePrepare Done ++++++")
}

// HandleP1a handles Pre_prepare message
func (p *srbft) HandlePre(m PrePrepare) {
	log.Debugf("<--------------------HandlePre------------------>")

	log.Debugf(" Sender  %v ", m.ID)

	log.Debugf(" m.Slot  %v ", m.Slot)

	if m.Ballot > p.ballot {
		log.Debugf("m.Ballot > p.ballot")
		p.ballot = m.Ballot
		p.view = m.View
	}

	_, ok := p.log[m.Slot]
	if !ok {
		log.Debugf("Create a log")
		p.log[m.Slot] = &entry{
			ballot:    p.ballot,
			view:      p.view,
			command:   m.Command,
			commit:    false,
			active:    false,
			Leader:    false,
			request:   &m.Request,
			timestamp: time.Now(),
			Digest:    GetMD5Hash(&m.Request),
			Q1:        PaxiBFT.NewQuorum(),
			Q2:        PaxiBFT.NewQuorum(),
			Q3:        PaxiBFT.NewQuorum(),
			Q4:        PaxiBFT.NewQuorum(),
		}
	}
	e := p.log[m.Slot]

	if e.request == nil {
		e.request = &m.Request
	}

	e.Digest = GetMD5Hash(&m.Request)
	for i, v := range e.Digest {
		if v != m.Digest[i] {
			return
		}
	}
	log.Debugf("m.Ballot=%v , p.ballot=%v, m.view=%v", m.Ballot, p.ballot, m.View)
	log.Debugf("at the prepare handling")
	p.Broadcast(Prepare{
		Ballot: p.ballot,
		ID:     p.ID(),
		View:   m.View,
		Slot:   m.Slot,
		Digest: m.Digest,
	})
	log.Debugf("++++++ HandlePre Done ++++++")
	p.HandlePrepare(Prepare{
		Ballot: p.ballot,
		ID:     p.ID(),
		View:   m.View,
		Slot:   m.Slot,
		Digest: m.Digest,
	})
}

func (p *srbft) LightPrePrepare(r *PaxiBFT.Request, slt int) {
	log.Debugf("<--------------------LightPrePrepare------------------>")
	s := GetMD5Hash(r)

	//log.Debugf("strugglerThreshold: %v, N: %v", p.strugglerThreshold, p.N.N())
	for id := p.strugglerThreshold; id <= 4; id++ {
		if id != p.ID().Node() {

			p.Send(PaxiBFT.NewID(1, id), LightPrePrepare{
				Ballot: p.ballot,
				ID:     p.ID(),
				View:   p.view,
				Slot:   slt,
				Digest: s,
			})
			log.Debugf("LightPrePrepare sent to %v", id)
		}
	}
}

func (p *srbft) HandleLightPre(m LightPrePrepare) {
	log.Debugf("<--------------------HandleLightPre------------------>")

	log.Debugf(" Sender  %v ", m.ID)

	log.Debugf(" m.Slot  %v ", m.Slot)

	if m.Ballot > p.ballot {
		log.Debugf("m.Ballot > p.ballot")
		p.ballot = m.Ballot
		p.view = m.View
	}

	_, ok := p.log[m.Slot]
	if !ok {
		log.Debugf("Create a log")
		p.log[m.Slot] = &entry{
			ballot:    p.ballot,
			view:      p.view,
			commit:    false,
			active:    false,
			Leader:    false,
			request:   nil,
			timestamp: time.Now(),
			Digest:    m.Digest,
			Q1:        PaxiBFT.NewQuorum(),
			Q2:        PaxiBFT.NewQuorum(),
			Q3:        PaxiBFT.NewQuorum(),
			Q4:        PaxiBFT.NewQuorum(),
		}
	}

	log.Debugf("m.Ballot=%v , p.ballot=%v, m.view=%v", m.Ballot, p.ballot, m.View)
	log.Debugf("++++++ HandleLightPre Done ++++++")

}

// HandlePrepare starts phase 2 HandlePrepare
func (p *srbft) HandlePrepare(m Prepare) {
	log.Debugf("<--------------------HandlePrepare------------------>")
	log.Debugf(" Sender  %v ", m.ID)
	log.Debugf("p.slot=%v", p.slot)
	log.Debugf("m.slot=%v", m.Slot)

	_, ok := p.log[m.Slot]

	if !ok {
		log.Debugf("we create a log")
		p.log[m.Slot] = &entry{
			ballot:    p.ballot,
			view:      p.view,
			command:   m.Command,
			commit:    false,
			active:    false,
			Leader:    false,
			request:   nil,
			timestamp: time.Now(),
			Digest:    m.Digest,
			Q1:        PaxiBFT.NewQuorum(),
			Q2:        PaxiBFT.NewQuorum(),
			Q3:        PaxiBFT.NewQuorum(),
			Q4:        PaxiBFT.NewQuorum(),
		}
	}
	e := p.log[m.Slot]
	e.Q1.ACK(m.ID)

	if e.Q1.PreparedMajority() {
		e.Q1.Reset()
		e.Pstatus = PREPARED
		p.Broadcast(Commit{
			Ballot: p.ballot,
			ID:     p.ID(),
			View:   p.view,
			Slot:   m.Slot,
			Digest: m.Digest,
		})
	} else if p.ID().Node() >= p.strugglerThreshold { // if this node is struggler and the vote is required
		if e.Q1.StrugglerQuorum() {
			log.Debugf("Struggler quorum Needed")
			p.Broadcast(Prepare{
				Ballot: p.ballot,
				ID:     p.ID(),
				View:   m.View,
				Slot:   m.Slot,
				Digest: m.Digest,
			})
			if e.Q1.PreparedMajority() {
				e.Pstatus = PREPARED
				p.Broadcast(Commit{
					Ballot: p.ballot,
					ID:     p.ID(),
					View:   p.view,
					Slot:   m.Slot,
					Digest: m.Digest,
				})
			}
		}
	}

	if e.Cstatus == COMMITTED && e.Pstatus == PREPARED && e.Rstatus == RECEIVED {
		e.commit = true
		p.exec()
	}
	log.Debugf("++++++ HandlePrepare Done ++++++")
}

// HandleCommit starts phase 3
func (p *srbft) HandleCommit(m Commit) {
	log.Debugf("<--------------------HandleCommit------------------>")
	log.Debugf(" Sender  %v ", m.ID)
	log.Debugf("m.slot=%v", m.Slot)
	log.Debugf("p.slot=%v", p.slot)
	if p.execute > m.Slot {
		log.Debugf("old message")
		return
	}
	_, exist := p.log[m.Slot]
	if !exist {
		log.Debugf("create a log")
		p.log[m.Slot] = &entry{
			ballot:    p.ballot,
			view:      p.view,
			command:   m.Command,
			commit:    false,
			active:    false,
			Leader:    false,
			request:   nil,
			timestamp: time.Now(),
			Digest:    m.Digest,
			Q1:        PaxiBFT.NewQuorum(),
			Q2:        PaxiBFT.NewQuorum(),
			Q3:        PaxiBFT.NewQuorum(),
			Q4:        PaxiBFT.NewQuorum(),
		}
	}
	e := p.log[m.Slot]
	e.Q2.ACK(m.ID)

	log.Debugf("Q2 size =%v", e.Q2.Size())
	if e.Q2.PreparedMajority() {
		e.Cstatus = COMMITTED
	}
	if (e.Q2.PreparedMajority() || e.Cstatus == COMMITTED) && e.Pstatus == PREPARED && e.Rstatus == RECEIVED {
		e.Q2.Reset()
		e.commit = true
		p.exec()
	}
	log.Debugf("********* Commit End *********** ")
}

func (p *srbft) exec() {
	log.Debugf("<--------------------exec()------------------>")
	for {
		log.Debugf("p.execute %v", p.execute)
		e, ok := p.log[p.execute]
		if !ok || !e.commit {
			log.Debugf("Break")
			break
		}
		value := p.Execute(e.command)
		log.Debugf("value=%v", value)

		reply := PaxiBFT.Reply{
			Command:    e.command,
			Value:      value,
			Properties: make(map[string]string),
		}

		if e.request != nil && e.Leader {
			log.Debugf(" ********* Primary Request ********* %v", *e.request)
			e.request.Reply(reply)
			log.Debugf("********* Reply Primary *********")
			e.request = nil
		} else {
			log.Debugf("********* Replica Request ********* ")
			log.Debugf("p.ID() =%v", p.ID())
			e.request.Reply(reply)
			e.request = nil
			log.Debugf("********* Reply Replicas *********")
		}
		// TODO clean up the log periodically
		delete(p.log, p.execute)
		p.execute++
	}
}

func (p *srbft) BroadcastToStrugglers(m interface{}) {
	for id := p.strugglerThreshold; id <= 4; id++ {
		if id != p.ID().Node() {
			p.Send(PaxiBFT.NewID(1, id), m)
		}
	}
}

func (p *srbft) BroadcastToFast(m interface{}) {
	for id := 1; id < p.strugglerThreshold; id++ {
		if id != p.ID().Node() {
			p.Send(PaxiBFT.NewID(1, id), m)
		}
	}
}
