rm ./TestBunServer.tar.gz

rm ./app
# --target=bun-linux-x64
bun build --compile --minify --sourcemap  ./index.ts --outfile app
# bun build --compile --minify --sourcemap  --target=bun-linux-x64 ./index.ts --outfile app


# tar -czvf TestBunServer.tar.gz ./app ./src/client.tsx

# sgridnext

# rm ./app
