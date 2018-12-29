#!/usr/bin/env bash

B="\033[0;1m"
G="\033[0;92m"
R="\033[0m"

echo -e "${B}${G}Let's make some hash${R}"

echo -e "\n${B}${G}  - cleanup environment${R}"
echo -e "remove files: small.hash, large.hash, large.txt"
docker exec -it hash-demo sh -c "rm -rf small.hash"
docker exec -it hash-demo sh -c "rm -rf large.hash"
docker exec -it hash-demo sh -c "rm -rf large.txt"

echo -e "\n${B}${G}  - make large file (concat small files)${R}"
for i in $(seq -f "%02g" 10)
do
    echo "  #> cat $i.txt >> large.txt"
    docker exec -it hash-demo sh -c "cat $i.txt >> large.txt"
done

echo -e "\n${B}${G}  - make hash of small files${R}"
for i in $(seq -f "%02g" 10)
do
    echo -e "  #> homo -file $i.txt | tee -a small.hash"
    docker exec -it hash-demo sh -c "homo -file $i.txt | tee -a small.hash"
done

echo -e "\n${B}${G}  - make hash of large${R}"
echo -e "  #> homo -file large.txt | homo -concat"
docker exec -it hash-demo sh -c 'homo -file large.txt | homo -concat'

echo -e "\n${B}${G}  - make hash of pieces${R}"
echo -e "  #> cat small.hash | homo -concat"
docker exec -it hash-demo sh -c 'cat small.hash | homo -concat '