#Usage: gomsgbus <NODENAME> <CMDPORT> <MYURL> <URL>...
url0=tcp://127.0.0.1:40890
url1=tcp://127.0.0.1:40891
url2=tcp://127.0.0.1:40892
url3=tcp://127.0.0.1:40893


./killem.sh
./gomsgbus node0 1234 $url0 $url1 $url2 & node0=$!
./gomsgbus node1 2345 $url1 $url2 $url3 & node1=$!
./gomsgbus node2 3456 $url2 $url3 & node2=$!
./gomsgbus node3 4567 $url3 $url0 & node3=$!

#sleep 5
#kill $node0 $node1 $node2 $node3
