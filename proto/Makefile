# TODO(roy): These are currently not used. Clean this up when moving proto-js and proto-py.
# proto-js needs cosmos protos that aren't depended on by vindax protos.
# Running the export command above will only get a subset of cosmos protos.
# For now, we will export cosmos protos separately. Note that this doesn't verify the digest.
# TODO(DEC-2027): Do this in a more reproducible way.
proto-export-proto-js:
	rm -rf ./.proto-export
	buf export --config ./buf.yaml --output ./.proto-export --exclude-imports
	buf export buf.build/cosmos/cosmos-sdk:v0.50.0 --output ./.proto-export
