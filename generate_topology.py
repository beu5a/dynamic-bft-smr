import xml.etree.ElementTree as ET
from xml.dom import minidom
import argparse


def prettify(elem):
    """Return a pretty-printed XML string for the Element."""
    rough_string = ET.tostring(elem, 'utf-8')
    reparsed = minidom.parseString(rough_string)
    return reparsed.toprettyxml(indent="    ")

def create_service(name, image, replicas, share, command, supervisor=False, port=None):
    service = ET.Element("service", name=name, image=image, replicas=replicas, share=share, command=command)
    if supervisor:
        service.set("supervisor", "true")
    if port:
        service.set("port", port)
    return service

def create_link(origin, dest, latency, upload, download, drop, network):
    link = ET.Element("link", origin=origin, dest=dest, latency=latency, upload=upload, download=download, drop=drop, network=network)
    return link

def create_schedule(name, time, action):
    schedule = ET.Element("schedule", name=name, time=time, action=action)
    return schedule

def generate_topology(num_servers, num_strugglers, normal_link_capacity, struggler_link_capacity):
    root = ET.Element("experiment", boot="kollaps:2.0")

    services = ET.SubElement(root, "services")
    services.append(create_service("dashboard", "kollaps/dashboard:1.0", "1", "true", "['server']", supervisor=True, port="8088"))

    for i in range(num_servers):
        services.append(create_service(f"s{i}", "kollaps/paxibft:1.0", "1", "false", f"['server','1.{i+1}']"))

    services.append(create_service("client", "kollaps/paxibft:1.0", "1", "false", "['client','1.1']"))

    bridges = ET.SubElement(root, "bridges")
    
    links = ET.SubElement(root, "links")
    
    network = "kollaps_network"
    
    for i in range(num_servers):
        for j in range(i + 1, num_servers):
            capacity = struggler_link_capacity if i >= (num_servers-num_strugglers) or j >= (num_servers-num_strugglers)else normal_link_capacity
            links.append(create_link(f"s{i}", f"s{j}", "1", capacity, capacity, "0", network))

    client_capacity = "200Mbps"
    for i in range(num_servers):
        links.append(create_link("client", f"s{i}", "1", client_capacity, client_capacity, "0", network))

    dynamic = ET.SubElement(root, "dynamic")
    
    for i in range(num_servers):
        dynamic.append(create_schedule(f"s{i}", "0.0", "join"))

    dynamic.append(create_schedule("client", "0.0", "join"))

    return prettify(root)

if __name__ == "__main__":
    num_servers = 4
    num_strugglers = 2
    normal_link_capacity = "10Kbps"
    struggler_link_capacity = "1Kbps"
    client_capacity = "1000Mbps"

    # read arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("--num_servers", type=int, default=num_servers)
    parser.add_argument("--num_strugglers", type=int, default=num_strugglers)
    parser.add_argument("--normal_link_capacity", type=str, default=normal_link_capacity)
    parser.add_argument("--struggler_link_capacity", type=str, default=struggler_link_capacity)
    args = parser.parse_args()
    # example: python generate_topology.py --num_servers 5 --num_strugglers 2 --normal_link_capacity "100Kbps" --struggler_link_capacity "10Kbps"


    
    topology_xml = generate_topology(args.num_servers, args.num_strugglers, args.normal_link_capacity, args.struggler_link_capacity)
    with open("topology.xml", "w") as f:
        f.write(topology_xml)
    print("Topology XML generated and saved to topology.xml")
