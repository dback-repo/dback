containers=`docker ps --format '{{.Names}}'`
echo $containers

for i in $containers; do
    #echo $i
    docker inspect -f '{{.Mounts}}' $i
done