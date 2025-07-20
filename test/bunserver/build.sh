rm ./BunSSRTestServer.tar.gz

rm ./app
rm -rf ./dist
# --target=bun-linux-x64
# bun build --compile --minify --sourcemap  ./index.ts --outfile app
bun build --compile --minify --sourcemap  --target=bun-linux-x64 ./index.ts --outfile app

mkdir -p ./dist/src/client
cp ./app ./dist/app
cp ./src/client.tsx ./dist/src/client.tsx
cp -r ./src/client/* ./dist/src/client/
cp ./package.json ./dist/package.json
cp ./sgridnext.release ./dist/sgridnext.release

cd ./dist
bun i --production

tar -czvf BunSSRTestServer.tar.gz ./*

sgridnext

# rm ./app
