# Start the JSON file
config_file="config.json"

cat <<EOF > $config_file
{
    "address": {
EOF

# Append TCP addresses
grep 'tcp' hosts.temp | sed 's/$/,/' | sed '$ s/,$//' >> $config_file

cat <<EOF >> $config_file
    },
    "http_address": {
EOF

# Append HTTP addresses
grep 'http' hosts.temp | sed 's/$/,/' | sed '$ s/,$//' >> $config_file
echo '}}' >> "$config_file"

# close bracket


set -a
source /.env
set +a




# Create a new JSON object with the passed arguments
new_json=$(jq -n \
    --arg POLICY "$POLICY" \
    --arg THRESHOLD "$THRESHOLD" \
    --arg THRIFTY "$THRIFTY" \
    --arg CHAN_BUFFER_SIZE "$CHAN_BUFFER_SIZE" \
    --arg BUFFER_SIZE "$BUFFER_SIZE" \
    --arg MULTIVERSION "$MULTIVERSION" \
    --arg T "$T" \
    --arg N "$N" \
    --arg K "$K" \
    --arg W "$W" \
    --arg THROTTLE "$THROTTLE" \
    --arg CONCURRENCY "$CONCURRENCY" \
    --arg DISTRIBUTION "$DISTRIBUTION" \
    --arg LINEARIZABILITY_CHECK "$LINEARIZABILITY_CHECK" \
    --arg CONFLICTS "$CONFLICTS" \
    --arg MIN "$MIN" \
    --arg MU "$MU" \
    --arg SIGMA "$SIGMA" \
    --arg MOVE "$MOVE" \
    --arg SPEED "$SPEED" \
    --arg ZIPFIAN_S "$ZIPFIAN_S" \
    --arg ZIPFIAN_V "$ZIPFIAN_V" \
    --arg LAMBDA "$LAMBDA" \
    '{
        "policy": $POLICY,
        "threshold": ($THRESHOLD | tonumber),
        "thrifty": ($THRIFTY  | test("true")),
        "chan_buffer_size": ($CHAN_BUFFER_SIZE | tonumber),
        "buffer_size": ($BUFFER_SIZE | tonumber) ,
        "multiversion": ($MULTIVERSION | test("true")),
        "benchmark": {
            "T": ($T | tonumber),
            "N": ($N | tonumber),
            "K": ($K | tonumber),
            "W": ($W | tonumber),
            "Throttle": ($THROTTLE | tonumber),
            "Concurrency": ($CONCURRENCY | tonumber),
            "Distribution": $DISTRIBUTION,
            "LinearizabilityCheck": ($LINEARIZABILITY_CHECK | test("true")),
            "Conflicts": ($CONFLICTS | tonumber),
            "Min": ($MIN | tonumber),
            "Mu": ($MU | tonumber),
            "Sigma": ($SIGMA | tonumber),
            "Move": ($MOVE | test("true")),
            "Speed": ($SPEED | tonumber),
            "Zipfian_s": ($ZIPFIAN_S | tonumber),
            "Zipfian_v": ($ZIPFIAN_V | tonumber),
            "Lambda": ($LAMBDA | tonumber)
        }
    }')


# Check if config.json exists and has size greater than zero
if [ -s "$config_file" ]; then
    # Merge new JSON data into the existing JSON file
    jq -s '.[0] * .[1]' "$config_file" <(echo "$new_json") > tmp.$$ && mv tmp.$$ "$config_file"
else
    # Create new config.json if not present
    echo "$new_json" > "$config_file"
fi

echo "config.json updated or created"
