# gnalloy linux full parity matrix

Linux 完整对标矩阵：TCP、UDP、HTTP/1.1、HTTP/2、HTTP/3、HTTPS TLS 1.1/1.2/1.3 与纯 QUIC stream。HTTP/2 over TLS 跳过 TLS 1.1，HTTP/3/QUIC 固定 TLS 1.3；gnet、fasthttp、netpoll 仅执行各自真实支持的等价协议。

## Machine

| Field | Value |
| --- | --- |
| hostname | xigexb-dev2 |
| os | linux |
| arch | amd64 |
| cpus | 8 |
| go | go1.27.0 |
| ips | 172.16.8.171,172.30.0.1,172.19.0.1,192.168.16.1,192.168.32.1,172.23.0.1,172.21.0.1,172.25.0.1,172.17.0.1,10.255.255.1,10.128.0.1,10.128.16.1,10.128.32.1,10.128.48.1,10.128.64.1,10.128.80.1,10.128.96.1,10.128.112.1,10.128.128.1,10.128.144.1,10.128.160.1,fe80::7894:c91d:a414:de74,fe80::e405:35ff:feb9:2c8c,fe80::189f:87ff:fe0b:134f,fe80::44c5:f2ff:fef2:c10b,fe80::c87b:4eff:fe33:d286,fe80::a488:efff:fec0:4995,fe80::34ea:baff:fe74:a293,fe80::a8e6:a7ff:fe1b:92dc,fe80::d0f4:ecff:fe3d:d45f,fe80::d00b:8cff:fe9e:2817,fe80::f06e:f1ff:fe2e:c7ab,fe80::109b:bff:fe67:d473,fe80::8872:c4ff:fe9c:3c86,fe80::c8ef:1aff:fe7e:b235,fe80::8883:a2ff:fef9:28e8,fe80::e8d7:9fff:fefc:13db,fe80::2cf5:42ff:fe0e:83fd,fe80::90f9:9cff:fe99:a8a8,fe80::430:d0ff:feba:4364,fe80::646c:c7ff:fe99:f64c,fe80::de:20ff:fe8e:ef5d,fe80::40c1:fff:fe68:8139,fe80::42c:67ff:fe57:c074,fe80::ec96:ebff:fe6c:831d,fe80::f06f:deff:fe6c:7c74,fe80::5c46:8bff:fefe:b31f,fe80::ac44:42ff:fe6e:cc97,fe80::84aa:b0ff:fe1a:c6a,fe80::cb4:89ff:fe44:f8a7,fe80::74c9:9bff:fe1b:5fb8,fe80::8e0:f6ff:fe4b:ea06,fe80::6c30:d5ff:fe92:e7d7,fe80::c09d:f0ff:fe43:61bf,fe80::48ab:8eff:fedd:9fbe,fe80::70a3:24ff:feed:2b42,fe80::1cdd:84ff:fe15:1624,fe80::64c2:12ff:fe88:476a,fe80::2c3b:1aff:fed7:3b9c,fe80::468:ddff:fe87:3aeb,fe80::d02b:ff:fe75:30da,fe80::2c07:a0ff:fe1a:e3e8,fe80::8861:55ff:fe92:d3a8 |

## Summary

| Scenario | Framework | Protocol | ALPN | Backend | Loops | Total | Errors | Throughput ops/s | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy epoll tcp-echo 64B | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 195715.48 | 2098731 | 17702912 | 0 | 1.635026519s |
| gnalloy epoll tcp-echo 64B | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 158226.41 | 2503377 | 17567744 | 0 | 2.022418385s |
| gnalloy epoll tcp-echo 64B | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 184763.77 | 2213392 | 19808256 | 0 | 1.731941271s |
| netty epoll tcp-echo 64B | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 140364.37 | 3253290 | 174456832 | 0 | 2.279780891s |
| netty epoll tcp-echo 64B | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 174733.68 | 2323962 | 184086528 | 0 | 1.831358471s |
| netty epoll tcp-echo 64B | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 184024.35 | 1790865 | 187568128 | 0 | 1.73890027s |
| gnet poller tcp-echo 64B | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 159725.52 | 1915505 | 16527360 | 0 | 2.003436944s |
| gnet poller tcp-echo 64B | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 163247.79 | 1871223 | 16490496 | 0 | 1.960210345s |
| gnet poller tcp-echo 64B | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 159849.23 | 2007350 | 16637952 | 0 | 2.001886435s |
| netpoll poller tcp-echo 64B | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 227932.85 | 813925 | 19472384 | 10 | 1.403922257s |
| netpoll poller tcp-echo 64B | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 227643.66 | 591290 | 19181568 | 11 | 1.405705728s |
| netpoll poller tcp-echo 64B | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 227590.05 | 602507 | 17031168 | 11 | 1.406036888s |
| gnalloy epoll tcp-echo 1KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 199501.91 | 2204752 | 17604608 | 0 | 1.60399466s |
| gnalloy epoll tcp-echo 1KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 210033.86 | 2038075 | 17707008 | 0 | 1.523563855s |
| gnalloy epoll tcp-echo 1KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 240689.41 | 1257934 | 17440768 | 0 | 1.329514256s |
| netty epoll tcp-echo 1KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 166390.59 | 2779887 | 190599168 | 0 | 1.923185676s |
| netty epoll tcp-echo 1KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 149751.37 | 4140894 | 184356864 | 0 | 2.136875283s |
| netty epoll tcp-echo 1KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 155713.63 | 2667913 | 183382016 | 0 | 2.055054523s |
| gnet poller tcp-echo 1KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 155118.81 | 1825926 | 16490496 | 0 | 2.062934838s |
| gnet poller tcp-echo 1KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 157284.97 | 1989227 | 16650240 | 0 | 2.034523732s |
| gnet poller tcp-echo 1KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 160912.16 | 1775060 | 14426112 | 0 | 1.988662575s |
| netpoll poller tcp-echo 1KiB | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 208064.37 | 1035964 | 21393408 | 9 | 1.537985602s |
| netpoll poller tcp-echo 1KiB | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 205228.19 | 1333689 | 21106688 | 9 | 1.559239952s |
| netpoll poller tcp-echo 1KiB | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 222839.01 | 697326 | 17567744 | 9 | 1.436014297s |
| gnalloy epoll tcp-echo 16KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=16384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 160735.04 | 1924686 | 20230144 | 1 | 1.990854045s |
| gnalloy epoll tcp-echo 16KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=16384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 137732.87 | 2540461 | 20279296 | 1 | 2.323337931s |
| gnalloy epoll tcp-echo 16KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=16384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 149743.38 | 2224978 | 20307968 | 1 | 2.136989261s |
| netty epoll tcp-echo 16KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 127200.83 | 3116394 | 191815680 | 1 | 2.515706819s |
| netty epoll tcp-echo 16KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 115189.66 | 3554574 | 202924032 | 1 | 2.7780272s |
| netty epoll tcp-echo 16KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 89815.21 | 5016457 | 188547072 | 1 | 3.56287081s |
| gnet poller tcp-echo 16KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 116333.38 | 2806307 | 19206144 | 1 | 2.750715315s |
| gnet poller tcp-echo 16KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 132664.88 | 2251497 | 17113088 | 1 | 2.412092759s |
| gnet poller tcp-echo 16KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 320000 | 0 | 117960.79 | 2612989 | 19177472 | 1 | 2.712765829s |
| netpoll poller tcp-echo 16KiB | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 167505.02 | 774620 | 26025984 | 6 | 1.910390456s |
| netpoll poller tcp-echo 16KiB | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 159550.08 | 1135772 | 24215552 | 7 | 2.005639897s |
| netpoll poller tcp-echo 16KiB | netpoll | tcp-echo |  | poller | - | 320000 | 0 | 148872.84 | 1494665 | 24371200 | 6 | 2.149485359s |
| gnalloy epoll udp-echo 128B | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 111547.76 | 856062 | 18948096 | 51 | 2.868726387s |
| gnalloy epoll udp-echo 128B | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 99938.52 | 2869700 | 18640896 | 50 | 3.201968543s |
| gnalloy epoll udp-echo 128B | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 111718.52 | 834035 | 19148800 | 50 | 2.864341539s |
| netty epoll udp-echo 128B | netty | udp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 95052.13 | 903855 | 184332288 | 1 | 3.36657383s |
| netty epoll udp-echo 128B | netty | udp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 95573.97 | 877503 | 185024512 | 1 | 3.348191954s |
| netty epoll udp-echo 128B | netty | udp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 95856.29 | 900319 | 191569920 | 1 | 3.338330821s |
| gnet poller udp-echo 128B | gnet | udp-echo |  | poller | eventLoops=8 | 320000 | 0 | 176659.05 | 1538254 | 17375232 | 58 | 1.811398831s |
| gnet poller udp-echo 128B | gnet | udp-echo |  | poller | eventLoops=8 | 320000 | 0 | 165374.95 | 1998772 | 18132992 | 58 | 1.934996819s |
| gnet poller udp-echo 128B | gnet | udp-echo |  | poller | eventLoops=8 | 320000 | 0 | 170423.68 | 1520535 | 18169856 | 58 | 1.877673322s |
| gnalloy epoll udp-echo 1KiB | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 97810.54 | 1038865 | 19320832 | 48 | 3.271631125s |
| gnalloy epoll udp-echo 1KiB | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 106762.26 | 919077 | 21143552 | 48 | 2.997313711s |
| gnalloy epoll udp-echo 1KiB | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 98943.05 | 946924 | 18518016 | 48 | 3.234183855s |
| netty epoll udp-echo 1KiB | netty | udp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 90769.55 | 1475117 | 190017536 | 1 | 3.525411361s |
| netty epoll udp-echo 1KiB | netty | udp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 91599.83 | 986485 | 190365696 | 1 | 3.493456251s |
| netty epoll udp-echo 1KiB | netty | udp-echo |  | epoll | eventLoops=8 | 320000 | 0 | 92161.73 | 963504 | 195325952 | 1 | 3.472157208s |
| gnet poller udp-echo 1KiB | gnet | udp-echo |  | poller | eventLoops=8 | 320000 | 0 | 154395.3 | 2133582 | 18018304 | 61 | 2.072601933s |
| gnet poller udp-echo 1KiB | gnet | udp-echo |  | poller | eventLoops=8 | 320000 | 0 | 168307.35 | 1553677 | 18423808 | 60 | 1.901283518s |
| gnet poller udp-echo 1KiB | gnet | udp-echo |  | poller | eventLoops=8 | 320000 | 0 | 165959.31 | 1465318 | 18612224 | 60 | 1.928183508s |
| gnalloy epoll http1 128B | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 214698 | 1816515 | 22425600 | 2 | 1.490465671s |
| gnalloy epoll http1 128B | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 229384.51 | 1519743 | 22609920 | 2 | 1.395037556s |
| gnalloy epoll http1 128B | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 232896.85 | 1340394 | 22482944 | 2 | 1.373998814s |
| netty epoll http1 128B | netty | http1 |  | epoll | eventLoops=8 | 320000 | 0 | 107422.28 | 3945852 | 372654080 | 7 | 2.978897936s |
| netty epoll http1 128B | netty | http1 |  | epoll | eventLoops=8 | 320000 | 0 | 111199.45 | 3344587 | 418213888 | 6 | 2.877712045s |
| netty epoll http1 128B | netty | http1 |  | epoll | eventLoops=8 | 320000 | 0 | 121872.99 | 3582067 | 322428928 | 7 | 2.625684263s |
| gnet poller http1 128B | gnet | http1 |  | poller | eventLoops=8 | 320000 | 0 | 152682.43 | 1832194 | 21409792 | 2 | 2.095853444s |
| gnet poller http1 128B | gnet | http1 |  | poller | eventLoops=8 | 320000 | 0 | 147109.96 | 2373167 | 21635072 | 2 | 2.17524356s |
| gnet poller http1 128B | gnet | http1 |  | poller | eventLoops=8 | 320000 | 0 | 152701.08 | 1545825 | 17272832 | 1 | 2.095597434s |
| fasthttp http1 128B | fasthttp | http1 |  | net | - | 320000 | 0 | 175212.4 | 1947699 | 19243008 | 2 | 1.826354799s |
| fasthttp http1 128B | fasthttp | http1 |  | net | - | 320000 | 0 | 152210.02 | 2605204 | 19161088 | 1 | 2.102358265s |
| fasthttp http1 128B | fasthttp | http1 |  | net | - | 320000 | 0 | 173751.22 | 1693118 | 19116032 | 1 | 1.841713693s |
| netpoll poller http1 128B | netpoll | http1 |  | poller | - | 320000 | 0 | 227589.98 | 656706 | 22401024 | 15 | 1.406037282s |
| netpoll poller http1 128B | netpoll | http1 |  | poller | - | 320000 | 0 | 205297.06 | 1064011 | 21094400 | 14 | 1.558716887s |
| netpoll poller http1 128B | netpoll | http1 |  | poller | - | 320000 | 0 | 206996.04 | 994058 | 21434368 | 14 | 1.545923263s |
| gnalloy epoll http1 1KiB | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 230034.58 | 1587214 | 22286336 | 2 | 1.391095188s |
| gnalloy epoll http1 1KiB | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 230302.39 | 1368974 | 20144128 | 2 | 1.389477563s |
| gnalloy epoll http1 1KiB | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 196377.24 | 2227895 | 20353024 | 2 | 1.629516736s |
| netty epoll http1 1KiB | netty | http1 |  | epoll | eventLoops=8 | 320000 | 0 | 99803.35 | 4267286 | 325349376 | 7 | 3.206305229s |
| netty epoll http1 1KiB | netty | http1 |  | epoll | eventLoops=8 | 320000 | 0 | 92624.89 | 4964408 | 396431360 | 7 | 3.454794901s |
| netty epoll http1 1KiB | netty | http1 |  | epoll | eventLoops=8 | 320000 | 0 | 108438.53 | 3938645 | 379572224 | 6 | 2.950980624s |
| gnet poller http1 1KiB | gnet | http1 |  | poller | eventLoops=8 | 320000 | 0 | 148793.4 | 1823550 | 21479424 | 2 | 2.150632975s |
| gnet poller http1 1KiB | gnet | http1 |  | poller | eventLoops=8 | 320000 | 0 | 149768.37 | 2135331 | 21417984 | 1 | 2.13663269s |
| gnet poller http1 1KiB | gnet | http1 |  | poller | eventLoops=8 | 320000 | 0 | 150384.52 | 1659810 | 19476480 | 2 | 2.127878541s |
| fasthttp http1 1KiB | fasthttp | http1 |  | net | - | 320000 | 0 | 175041.37 | 1442534 | 21114880 | 1 | 1.828139208s |
| fasthttp http1 1KiB | fasthttp | http1 |  | net | - | 320000 | 0 | 148220.33 | 2829956 | 19140608 | 1 | 2.158948048s |
| fasthttp http1 1KiB | fasthttp | http1 |  | net | - | 320000 | 0 | 173618.58 | 1932874 | 19124224 | 2 | 1.843120752s |
| netpoll poller http1 1KiB | netpoll | http1 |  | poller | - | 320000 | 0 | 204169.58 | 1177092 | 24276992 | 13 | 1.567324571s |
| netpoll poller http1 1KiB | netpoll | http1 |  | poller | - | 320000 | 0 | 205662.31 | 1038564 | 23552000 | 13 | 1.555948695s |
| netpoll poller http1 1KiB | netpoll | http1 |  | poller | - | 320000 | 0 | 213415.95 | 921344 | 21184512 | 13 | 1.499419316s |
| gnalloy epoll https1 tls11 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 79085.31 | 3350222 | 33824768 | 79 | 4.046263488s |
| gnalloy epoll https1 tls11 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 79930.44 | 3269961 | 34107392 | 78 | 4.003481011s |
| gnalloy epoll https1 tls11 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 90247.47 | 2384097 | 34328576 | 79 | 3.545805839s |
| netty epoll https1 tls11 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 57265.49 | 5808241 | 398970880 | 13 | 5.588008089s |
| netty epoll https1 tls11 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 58411.99 | 6444794 | 481230848 | 12 | 5.478327651s |
| netty epoll https1 tls11 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 70047.47 | 4793318 | 408506368 | 13 | 4.568330367s |
| fasthttp https1 tls11 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 108258.84 | 2407510 | 22130688 | 18 | 2.95587866s |
| fasthttp https1 tls11 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 111350.99 | 2173329 | 23277568 | 19 | 2.873795676s |
| fasthttp https1 tls11 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 121235.98 | 1377070 | 22765568 | 20 | 2.639480501s |
| gnalloy epoll https1 tls11 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 80765.56 | 2772862 | 34344960 | 77 | 3.962084652s |
| gnalloy epoll https1 tls11 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 67860.32 | 4384733 | 34045952 | 75 | 4.715568535s |
| gnalloy epoll https1 tls11 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 80862.4 | 2483285 | 34226176 | 76 | 3.957339665s |
| netty epoll https1 tls11 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 55418.3 | 6388434 | 424497152 | 18 | 5.774266216s |
| netty epoll https1 tls11 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 56772.24 | 6314190 | 472879104 | 16 | 5.636557434s |
| netty epoll https1 tls11 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 59380.59 | 5426107 | 407826432 | 18 | 5.388966547s |
| fasthttp https1 tls11 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 86802.82 | 3747313 | 23191552 | 18 | 3.686516053s |
| fasthttp https1 tls11 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 106342.97 | 1647843 | 20987904 | 19 | 3.009131743s |
| fasthttp https1 tls11 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 89683.74 | 3580365 | 24707072 | 18 | 3.568093703s |
| gnalloy epoll https1 tls12 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 107850.14 | 2676288 | 33828864 | 73 | 2.967080049s |
| gnalloy epoll https1 tls12 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 129077.46 | 1461935 | 31649792 | 69 | 2.479131545s |
| gnalloy epoll https1 tls12 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 128757.17 | 1351212 | 33808384 | 67 | 2.48529852s |
| netty epoll https1 tls12 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 64245.77 | 5648456 | 475123712 | 17 | 4.980872844s |
| netty epoll https1 tls12 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 65881.36 | 6198546 | 454008832 | 17 | 4.85721631s |
| netty epoll https1 tls12 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 67317.23 | 6735373 | 605855744 | 14 | 4.753611882s |
| fasthttp https1 tls12 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 134742.56 | 2716256 | 22441984 | 3 | 2.374899296s |
| fasthttp https1 tls12 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 160357.95 | 1622029 | 22396928 | 4 | 1.995535596s |
| fasthttp https1 tls12 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 184897.32 | 887788 | 22245376 | 3 | 1.730690289s |
| gnalloy epoll https1 tls12 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 105829.19 | 2971700 | 34320384 | 68 | 3.023740431s |
| gnalloy epoll https1 tls12 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 104103.63 | 2613490 | 34476032 | 67 | 3.073860152s |
| gnalloy epoll https1 tls12 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 104872.87 | 2869251 | 33828864 | 68 | 3.051313597s |
| netty epoll https1 tls12 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 65785.66 | 6048609 | 469520384 | 17 | 4.864282077s |
| netty epoll https1 tls12 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 61472.63 | 6650647 | 557064192 | 16 | 5.205568357s |
| netty epoll https1 tls12 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 63092.69 | 5931033 | 556355584 | 16 | 5.071903118s |
| fasthttp https1 tls12 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 156260.76 | 1563392 | 20262912 | 3 | 2.047859036s |
| fasthttp https1 tls12 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 154474.08 | 2115709 | 22208512 | 3 | 2.071544905s |
| fasthttp https1 tls12 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 178867.65 | 919472 | 21839872 | 4 | 1.789032322s |
| gnalloy epoll https1 tls13 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 110555.13 | 3067208 | 33853440 | 61 | 2.894483426s |
| gnalloy epoll https1 tls13 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 108278.98 | 2762250 | 33980416 | 62 | 2.955328841s |
| gnalloy epoll https1 tls13 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 95126.75 | 3464141 | 33755136 | 62 | 3.363933002s |
| netty epoll https1 tls13 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 54217.07 | 7709098 | 638709760 | 14 | 5.902199742s |
| netty epoll https1 tls13 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 67738.03 | 5082296 | 638586880 | 15 | 4.724081562s |
| netty epoll https1 tls13 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 57739.94 | 6375818 | 495542272 | 17 | 5.542090962s |
| fasthttp https1 tls13 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 154036.7 | 1914047 | 24629248 | 9 | 2.077427067s |
| fasthttp https1 tls13 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 163333.99 | 1365796 | 22573056 | 7 | 1.959175808s |
| fasthttp https1 tls13 128B | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 182604.77 | 942106 | 24563712 | 8 | 1.75241862s |
| gnalloy epoll https1 tls13 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 115500.34 | 2579287 | 34033664 | 60 | 2.770554536s |
| gnalloy epoll https1 tls13 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 87113.36 | 4412951 | 32002048 | 59 | 3.673374399s |
| gnalloy epoll https1 tls13 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 104666.99 | 2652940 | 34029568 | 60 | 3.057315433s |
| netty epoll https1 tls13 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 53374.76 | 8248798 | 489738240 | 19 | 5.995343156s |
| netty epoll https1 tls13 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 47922.01 | 6595478 | 572129280 | 17 | 6.677515907s |
| netty epoll https1 tls13 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 320000 | 0 | 55306.52 | 6915962 | 642633728 | 15 | 5.785935793s |
| fasthttp https1 tls13 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 174828.25 | 989190 | 24592384 | 7 | 1.830367854s |
| fasthttp https1 tls13 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 176641.73 | 949962 | 22523904 | 7 | 1.811576454s |
| fasthttp https1 tls13 1KiB | fasthttp | https1 | http/1.1 | net | - | 320000 | 0 | 161524.15 | 1235495 | 22499328 | 7 | 1.981127928s |
| gnalloy epoll http2 128B | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 111205.26 | 1710540 | 22220800 | 118 | 2.877561691s |
| gnalloy epoll http2 128B | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 97102.02 | 2717045 | 23523328 | 116 | 3.29550295s |
| gnalloy epoll http2 128B | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 104052.77 | 2248865 | 24104960 | 115 | 3.075362497s |
| netty epoll http2 128B | netty | http2 |  | epoll | eventLoops=8 | 320000 | 0 | 64686.4 | 5944494 | 318230528 | 5 | 4.946943855s |
| netty epoll http2 128B | netty | http2 |  | epoll | eventLoops=8 | 320000 | 0 | 57423.77 | 7161116 | 322060288 | 5 | 5.572604866s |
| netty epoll http2 128B | netty | http2 |  | epoll | eventLoops=8 | 320000 | 0 | 72302.89 | 5161723 | 314843136 | 5 | 4.425825611s |
| gnalloy epoll http2 1KiB | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 93451.09 | 2883792 | 23633920 | 109 | 3.424251081s |
| gnalloy epoll http2 1KiB | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 110315.93 | 1723684 | 23724032 | 112 | 2.900759653s |
| gnalloy epoll http2 1KiB | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 96748.74 | 2746669 | 23928832 | 110 | 3.307536665s |
| netty epoll http2 1KiB | netty | http2 |  | epoll | eventLoops=8 | 320000 | 0 | 65525.81 | 5452168 | 331104256 | 5 | 4.883571776s |
| netty epoll http2 1KiB | netty | http2 |  | epoll | eventLoops=8 | 320000 | 0 | 66206.86 | 5242324 | 339124224 | 5 | 4.833335701s |
| netty epoll http2 1KiB | netty | http2 |  | epoll | eventLoops=8 | 320000 | 0 | 71927.33 | 4590663 | 307015680 | 5 | 4.448934785s |
| gnalloy epoll https2 tls12 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 60754.82 | 3722814 | 34066432 | 158 | 5.267071791s |
| gnalloy epoll https2 tls12 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 66598.24 | 2289368 | 34111488 | 160 | 4.804931579s |
| gnalloy epoll https2 tls12 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 59932.75 | 3751254 | 33533952 | 162 | 5.339317711s |
| netty epoll https2 tls12 128B | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 47847.08 | 6998574 | 481566720 | 22 | 6.687973352s |
| netty epoll https2 tls12 128B | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 43463.47 | 7597433 | 515072000 | 22 | 7.362504917s |
| netty epoll https2 tls12 128B | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 42111.28 | 6770355 | 620199936 | 17 | 7.598915011s |
| gnalloy epoll https2 tls12 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 57452.63 | 4183646 | 34041856 | 153 | 5.569805781s |
| gnalloy epoll https2 tls12 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 65847.2 | 2354118 | 33648640 | 154 | 4.859735623s |
| gnalloy epoll https2 tls12 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 58709.64 | 3762517 | 35962880 | 155 | 5.450552487s |
| netty epoll https2 tls12 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 45998.12 | 6364597 | 542945280 | 20 | 6.956806426s |
| netty epoll https2 tls12 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 43682.08 | 8228167 | 419569664 | 25 | 7.325657767s |
| netty epoll https2 tls12 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 43235.26 | 7733796 | 477458432 | 23 | 7.401366541s |
| gnalloy epoll https2 tls13 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 60696.97 | 3574613 | 34082816 | 156 | 5.272091772s |
| gnalloy epoll https2 tls13 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 65468.86 | 2826053 | 33566720 | 156 | 4.887820103s |
| gnalloy epoll https2 tls13 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 59244.87 | 3819385 | 34009088 | 155 | 5.401311761s |
| netty epoll https2 tls13 128B | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 38659.62 | 7778366 | 519393280 | 21 | 8.277371382s |
| netty epoll https2 tls13 128B | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 40709.22 | 8194155 | 664047616 | 18 | 7.860627915s |
| netty epoll https2 tls13 128B | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 38925.35 | 9108579 | 519905280 | 23 | 8.220864797s |
| gnalloy epoll https2 tls13 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 60407.78 | 3665159 | 33906688 | 148 | 5.297330579s |
| gnalloy epoll https2 tls13 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 56416.08 | 4399558 | 33947648 | 148 | 5.672141844s |
| gnalloy epoll https2 tls13 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 66457.64 | 2418169 | 33804288 | 151 | 4.8150972s |
| netty epoll https2 tls13 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 41880.39 | 8981866 | 587468800 | 20 | 7.640807118s |
| netty epoll https2 tls13 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 38143.19 | 7169950 | 513060864 | 24 | 8.389439893s |
| netty epoll https2 tls13 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 320000 | 0 | 36408.43 | 9003949 | 507047936 | 24 | 8.789172746s |
| gnalloy rfc9000 http3 tls13 128B | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 28225.85 | 4456117 | 1506222080 | 45 | 11.337125677s |
| gnalloy rfc9000 http3 tls13 128B | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 29463.05 | 4137028 | 1416114176 | 47 | 10.861063119s |
| gnalloy rfc9000 http3 tls13 128B | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 30093.25 | 3878332 | 1459339264 | 46 | 10.633613046s |
| netty epoll http3 tls13 128B | netty | http3 | h3 | epoll | eventLoops=8 | 320000 | 0 | 12973.94 | 10140036 | 423370752 | 20 | 24.664824847s |
| netty epoll http3 tls13 128B | netty | http3 | h3 | epoll | eventLoops=8 | 320000 | 0 | 13022.99 | 9844845 | 412467200 | 19 | 24.57192803s |
| netty epoll http3 tls13 128B | netty | http3 | h3 | epoll | eventLoops=8 | 320000 | 0 | 12768.43 | 10272673 | 412934144 | 20 | 25.061816195s |
| gnalloy rfc9000 http3 tls13 1KiB | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 28701.25 | 4500126 | 1343291392 | 45 | 11.149341652s |
| gnalloy rfc9000 http3 tls13 1KiB | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 30828.96 | 4291853 | 1432805376 | 45 | 10.37984988s |
| gnalloy rfc9000 http3 tls13 1KiB | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 30794.94 | 3848737 | 1471246336 | 44 | 10.391318812s |
| netty epoll http3 tls13 1KiB | netty | http3 | h3 | epoll | eventLoops=8 | 320000 | 0 | 12951.3 | 9913636 | 421994496 | 23 | 24.707941247s |
| netty epoll http3 tls13 1KiB | netty | http3 | h3 | epoll | eventLoops=8 | 320000 | 0 | 12795.95 | 9637870 | 425512960 | 23 | 25.007919644s |
| netty epoll http3 tls13 1KiB | netty | http3 | h3 | epoll | eventLoops=8 | 320000 | 0 | 12704.57 | 9605414 | 401076224 | 23 | 25.18779342s |
| gnalloy rfc9000 quic-stream tls13 128B | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 43161.9 | 5114955 | 37687296 | 232 | 7.413946059s |
| gnalloy rfc9000 quic-stream tls13 128B | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 45273.96 | 4756887 | 39854080 | 235 | 7.068079905s |
| gnalloy rfc9000 quic-stream tls13 128B | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 47532.81 | 4027335 | 37736448 | 235 | 6.732191305s |
| netty epoll quic-stream tls13 128B | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 320000 | 0 | 562.87 | 385753952 | 1238790144 | 31 | 9m28.510315805s |
| netty epoll quic-stream tls13 128B | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 320000 | 0 | 561.38 | 380746959 | 1247531008 | 31 | 9m30.0255379s |
| netty epoll quic-stream tls13 128B | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 320000 | 0 | 559.41 | 345269871 | 1244352512 | 31 | 9m32.033336759s |
| gnalloy rfc9000 quic-stream tls13 1KiB | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 47293.47 | 4166433 | 40570880 | 281 | 6.766262135s |
| gnalloy rfc9000 quic-stream tls13 1KiB | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 47150.35 | 4004240 | 38965248 | 278 | 6.786800755s |
| gnalloy rfc9000 quic-stream tls13 1KiB | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 320000 | 0 | 45062.68 | 4482656 | 39247872 | 277 | 7.101220499s |
| netty epoll quic-stream tls13 1KiB | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 320000 | 0 | 675.81 | 510589424 | 1153892352 | 30 | 7m53.504538557s |
| netty epoll quic-stream tls13 1KiB | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 320000 | 0 | 664.04 | 545051977 | 1157861376 | 30 | 8m1.897604131s |
| netty epoll quic-stream tls13 1KiB | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 320000 | 0 | 666.36 | 545080776 | 1153978368 | 31 | 8m0.218558545s |

| Scenario | Framework | Protocol | ALPN | Backend | Loops | Samples | Throughput min | Throughput median | Throughput max | Throughput mean | Median ns/op | Median P99 latency ns | Max RSS bytes | GC count | Errors |
| --- | --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy epoll tcp-echo 64B | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 158226.41 | 184763.77 | 195715.48 | 179568.55 | 5412.316471875 | 2213392 | 19808256 | 0 | 0 |
| netty epoll tcp-echo 64B | netty | tcp-echo |  | epoll | eventLoops=8 | 3 | 140364.37 | 174733.68 | 184024.35 | 166374.13 | 5722.995221875 | 2323962 | 187568128 | 0 | 0 |
| gnet poller tcp-echo 64B | gnet | tcp-echo |  | poller | eventLoops=8 | 3 | 159725.52 | 159849.23 | 163247.79 | 160940.85 | 6255.895109375 | 1915505 | 16637952 | 0 | 0 |
| netpoll poller tcp-echo 64B | netpoll | tcp-echo |  | poller | - | 3 | 227590.05 | 227643.66 | 227932.85 | 227722.19 | 4392.8304 | 602507 | 19472384 | 32 | 0 |
| gnalloy epoll tcp-echo 1KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 199501.91 | 210033.86 | 240689.41 | 216741.73 | 4761.137046875 | 2038075 | 17707008 | 0 | 0 |
| netty epoll tcp-echo 1KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 3 | 149751.37 | 155713.63 | 166390.59 | 157285.2 | 6422.045384375 | 2779887 | 190599168 | 0 | 0 |
| gnet poller tcp-echo 1KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 3 | 155118.81 | 157284.97 | 160912.16 | 157771.98 | 6357.8866625 | 1825926 | 16650240 | 0 | 0 |
| netpoll poller tcp-echo 1KiB | netpoll | tcp-echo |  | poller | - | 3 | 205228.19 | 208064.37 | 222839.01 | 212043.86 | 4806.20500625 | 1035964 | 21393408 | 27 | 0 |
| gnalloy epoll tcp-echo 16KiB | gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=16384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 137732.87 | 149743.38 | 160735.04 | 149403.76 | 6678.091440625 | 2224978 | 20307968 | 3 | 0 |
| netty epoll tcp-echo 16KiB | netty | tcp-echo |  | epoll | eventLoops=8 | 3 | 89815.21 | 115189.66 | 127200.83 | 110735.23 | 8681.335 | 3554574 | 202924032 | 3 | 0 |
| gnet poller tcp-echo 16KiB | gnet | tcp-echo |  | poller | eventLoops=8 | 3 | 116333.38 | 117960.79 | 132664.88 | 122319.68 | 8477.393215625 | 2612989 | 19206144 | 3 | 0 |
| netpoll poller tcp-echo 16KiB | netpoll | tcp-echo |  | poller | - | 3 | 148872.84 | 159550.08 | 167505.02 | 158642.65 | 6267.624678125 | 1135772 | 26025984 | 19 | 0 |
| gnalloy epoll udp-echo 128B | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 99938.52 | 111547.76 | 111718.52 | 107734.93 | 8964.769959375 | 856062 | 19148800 | 151 | 0 |
| netty epoll udp-echo 128B | netty | udp-echo |  | epoll | eventLoops=8 | 3 | 95052.13 | 95573.97 | 95856.29 | 95494.13 | 10463.09985625 | 900319 | 191569920 | 3 | 0 |
| gnet poller udp-echo 128B | gnet | udp-echo |  | poller | eventLoops=8 | 3 | 165374.95 | 170423.68 | 176659.05 | 170819.23 | 5867.72913125 | 1538254 | 18169856 | 174 | 0 |
| gnalloy epoll udp-echo 1KiB | gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 97810.54 | 98943.05 | 106762.26 | 101171.95 | 10106.824546875 | 946924 | 21143552 | 144 | 0 |
| netty epoll udp-echo 1KiB | netty | udp-echo |  | epoll | eventLoops=8 | 3 | 90769.55 | 91599.83 | 92161.73 | 91510.37 | 10917.050784375 | 986485 | 195325952 | 3 | 0 |
| gnet poller udp-echo 1KiB | gnet | udp-echo |  | poller | eventLoops=8 | 3 | 154395.3 | 165959.31 | 168307.35 | 162887.32 | 6025.5734625 | 1553677 | 18612224 | 181 | 0 |
| gnalloy epoll http1 128B | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 214698 | 229384.51 | 232896.85 | 225659.79 | 4359.4923625 | 1519743 | 22609920 | 6 | 0 |
| netty epoll http1 128B | netty | http1 |  | epoll | eventLoops=8 | 3 | 107422.28 | 111199.45 | 121872.99 | 113498.24 | 8992.850140625 | 3582067 | 418213888 | 20 | 0 |
| gnet poller http1 128B | gnet | http1 |  | poller | eventLoops=8 | 3 | 147109.96 | 152682.43 | 152701.08 | 150831.16 | 6549.5420125 | 1832194 | 21635072 | 5 | 0 |
| fasthttp http1 128B | fasthttp | http1 |  | net | - | 3 | 152210.02 | 173751.22 | 175212.4 | 167057.88 | 5755.355290625 | 1947699 | 19243008 | 4 | 0 |
| netpoll poller http1 128B | netpoll | http1 |  | poller | - | 3 | 205297.06 | 206996.04 | 227589.98 | 213294.36 | 4831.010196875 | 994058 | 22401024 | 43 | 0 |
| gnalloy epoll http1 1KiB | gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 196377.24 | 230034.58 | 230302.39 | 218904.74 | 4347.1724625 | 1587214 | 22286336 | 6 | 0 |
| netty epoll http1 1KiB | netty | http1 |  | epoll | eventLoops=8 | 3 | 92624.89 | 99803.35 | 108438.53 | 100288.92 | 10019.703840625 | 4267286 | 396431360 | 20 | 0 |
| gnet poller http1 1KiB | gnet | http1 |  | poller | eventLoops=8 | 3 | 148793.4 | 149768.37 | 150384.52 | 149648.76 | 6676.97715625 | 1823550 | 21479424 | 5 | 0 |
| fasthttp http1 1KiB | fasthttp | http1 |  | net | - | 3 | 148220.33 | 173618.58 | 175041.37 | 165626.76 | 5759.75235 | 1932874 | 21114880 | 4 | 0 |
| netpoll poller http1 1KiB | netpoll | http1 |  | poller | - | 3 | 204169.58 | 205662.31 | 213415.95 | 207749.28 | 4862.339671875 | 1038564 | 24276992 | 39 | 0 |
| gnalloy epoll https1 tls11 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 79085.31 | 79930.44 | 90247.47 | 83087.74 | 12510.878159375 | 3269961 | 34328576 | 236 | 0 |
| netty epoll https1 tls11 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 3 | 57265.49 | 58411.99 | 70047.47 | 61908.32 | 17119.773909375 | 5808241 | 481230848 | 38 | 0 |
| fasthttp https1 tls11 128B | fasthttp | https1 | http/1.1 | net | - | 3 | 108258.84 | 111350.99 | 121235.98 | 113615.27 | 8980.6114875 | 2173329 | 23277568 | 57 | 0 |
| gnalloy epoll https1 tls11 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 67860.32 | 80765.56 | 80862.4 | 76496.09 | 12381.5145375 | 2772862 | 34344960 | 228 | 0 |
| netty epoll https1 tls11 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 3 | 55418.3 | 56772.24 | 59380.59 | 57190.38 | 17614.24198125 | 6314190 | 472879104 | 52 | 0 |
| fasthttp https1 tls11 1KiB | fasthttp | https1 | http/1.1 | net | - | 3 | 86802.82 | 89683.74 | 106342.97 | 94276.51 | 11150.292821875 | 3580365 | 24707072 | 55 | 0 |
| gnalloy epoll https1 tls12 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 107850.14 | 128757.17 | 129077.46 | 121894.92 | 7766.557875 | 1461935 | 33828864 | 209 | 0 |
| netty epoll https1 tls12 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 3 | 64245.77 | 65881.36 | 67317.23 | 65814.79 | 15178.80096875 | 6198546 | 605855744 | 48 | 0 |
| fasthttp https1 tls12 128B | fasthttp | https1 | http/1.1 | net | - | 3 | 134742.56 | 160357.95 | 184897.32 | 159999.28 | 6236.0487375 | 1622029 | 22441984 | 10 | 0 |
| gnalloy epoll https1 tls12 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 104103.63 | 104872.87 | 105829.19 | 104935.23 | 9535.354990625 | 2869251 | 34476032 | 203 | 0 |
| netty epoll https1 tls12 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 3 | 61472.63 | 63092.69 | 65785.66 | 63450.33 | 15849.69724375 | 6048609 | 557064192 | 49 | 0 |
| fasthttp https1 tls12 1KiB | fasthttp | https1 | http/1.1 | net | - | 3 | 154474.08 | 156260.76 | 178867.65 | 163200.83 | 6399.5594875 | 1563392 | 22208512 | 10 | 0 |
| gnalloy epoll https1 tls13 128B | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 95126.75 | 108278.98 | 110555.13 | 104653.62 | 9235.402628125 | 3067208 | 33980416 | 185 | 0 |
| netty epoll https1 tls13 128B | netty | https1 | http/1.1 | epoll | eventLoops=8 | 3 | 54217.07 | 57739.94 | 67738.03 | 59898.35 | 17319.03425625 | 6375818 | 638709760 | 46 | 0 |
| fasthttp https1 tls13 128B | fasthttp | https1 | http/1.1 | net | - | 3 | 154036.7 | 163333.99 | 182604.77 | 166658.49 | 6122.4244 | 1365796 | 24629248 | 24 | 0 |
| gnalloy epoll https1 tls13 1KiB | gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 87113.36 | 104666.99 | 115500.34 | 102426.9 | 9554.110728125 | 2652940 | 34033664 | 179 | 0 |
| netty epoll https1 tls13 1KiB | netty | https1 | http/1.1 | epoll | eventLoops=8 | 3 | 47922.01 | 53374.76 | 55306.52 | 52201.1 | 18735.4473625 | 6915962 | 642633728 | 51 | 0 |
| fasthttp https1 tls13 1KiB | fasthttp | https1 | http/1.1 | net | - | 3 | 161524.15 | 174828.25 | 176641.73 | 170998.04 | 5719.89954375 | 989190 | 24592384 | 21 | 0 |
| gnalloy epoll http2 128B | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 97102.02 | 104052.77 | 111205.26 | 104120.02 | 9610.507803125 | 2248865 | 24104960 | 349 | 0 |
| netty epoll http2 128B | netty | http2 |  | epoll | eventLoops=8 | 3 | 57423.77 | 64686.4 | 72302.89 | 64804.35 | 15459.199546875 | 5944494 | 322060288 | 15 | 0 |
| gnalloy epoll http2 1KiB | gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 93451.09 | 96748.74 | 110315.93 | 100171.92 | 10336.052078125 | 2746669 | 23928832 | 331 | 0 |
| netty epoll http2 1KiB | netty | http2 |  | epoll | eventLoops=8 | 3 | 65525.81 | 66206.86 | 71927.33 | 67886.67 | 15104.174065625 | 5242324 | 339124224 | 15 | 0 |
| gnalloy epoll https2 tls12 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 59932.75 | 60754.82 | 66598.24 | 62428.6 | 16459.599346875 | 3722814 | 34111488 | 480 | 0 |
| netty epoll https2 tls12 128B | netty | https2 | h2 | epoll | eventLoops=8 | 3 | 42111.28 | 43463.47 | 47847.08 | 44473.94 | 23007.827865625 | 6998574 | 620199936 | 61 | 0 |
| gnalloy epoll https2 tls12 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 57452.63 | 58709.64 | 65847.2 | 60669.82 | 17032.976521875 | 3762517 | 35962880 | 462 | 0 |
| netty epoll https2 tls12 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 3 | 43235.26 | 43682.08 | 45998.12 | 44305.15 | 22892.680521875 | 7733796 | 542945280 | 68 | 0 |
| gnalloy epoll https2 tls13 128B | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 59244.87 | 60696.97 | 65468.86 | 61803.57 | 16475.2867875 | 3574613 | 34082816 | 467 | 0 |
| netty epoll https2 tls13 128B | netty | https2 | h2 | epoll | eventLoops=8 | 3 | 38659.62 | 38925.35 | 40709.22 | 39431.4 | 25690.202490625 | 8194155 | 664047616 | 62 | 0 |
| gnalloy epoll https2 tls13 1KiB | gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 56416.08 | 60407.78 | 66457.64 | 61093.83 | 16554.158059375 | 3665159 | 33947648 | 447 | 0 |
| netty epoll https2 tls13 1KiB | netty | https2 | h2 | epoll | eventLoops=8 | 3 | 36408.43 | 38143.19 | 41880.39 | 38810.67 | 26216.999665625 | 8981866 | 587468800 | 68 | 0 |
| gnalloy rfc9000 http3 tls13 128B | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 28225.85 | 29463.05 | 30093.25 | 29260.72 | 33940.822246875 | 4137028 | 1506222080 | 138 | 0 |
| netty epoll http3 tls13 128B | netty | http3 | h3 | epoll | eventLoops=8 | 3 | 12768.43 | 12973.94 | 13022.99 | 12921.79 | 77077.577646875 | 10140036 | 423370752 | 59 | 0 |
| gnalloy rfc9000 http3 tls13 1KiB | gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 28701.25 | 30794.94 | 30828.96 | 30108.38 | 32472.8712875 | 4291853 | 1471246336 | 134 | 0 |
| netty epoll http3 tls13 1KiB | netty | http3 | h3 | epoll | eventLoops=8 | 3 | 12704.57 | 12795.95 | 12951.3 | 12817.27 | 78149.7488875 | 9637870 | 425512960 | 69 | 0 |
| gnalloy rfc9000 quic-stream tls13 128B | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 43161.9 | 45273.96 | 47532.81 | 45322.89 | 22087.749703125 | 4756887 | 39854080 | 702 | 0 |
| netty epoll quic-stream tls13 128B | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 3 | 559.41 | 561.38 | 562.87 | 561.22 | 1781329.8059375 | 380746959 | 1247531008 | 93 | 0 |
| gnalloy rfc9000 quic-stream tls13 1KiB | gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 3 | 45062.68 | 47150.35 | 47293.47 | 46502.17 | 21208.752359375 | 4166433 | 40570880 | 836 | 0 |
| netty epoll quic-stream tls13 1KiB | netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 3 | 664.04 | 666.36 | 675.81 | 668.74 | 1500682.995453125 | 545051977 | 1157861376 | 91 | 0 |

| Scenario | Framework | Protocol | Benchmark | ns/op | B/op | allocs/op |
| --- | --- | --- | --- | ---: | ---: | ---: |
| gnalloy epoll tcp-echo 64B | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 5109 | 0 | 0 |
| gnalloy epoll tcp-echo 64B | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 6320 | 0 | 0 |
| gnalloy epoll tcp-echo 64B | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 5412 | 0 | 0 |
| netty epoll tcp-echo 64B | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 7124 | 0 | 0 |
| netty epoll tcp-echo 64B | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 5723 | 0 | 0 |
| netty epoll tcp-echo 64B | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 5434 | 0 | 0 |
| gnet poller tcp-echo 64B | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 6261 | 0 | 0 |
| gnet poller tcp-echo 64B | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 6126 | 0 | 0 |
| gnet poller tcp-echo 64B | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 6256 | 0 | 0 |
| netpoll poller tcp-echo 64B | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 4387 | 0 | 0 |
| netpoll poller tcp-echo 64B | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 4393 | 0 | 0 |
| netpoll poller tcp-echo 64B | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 4394 | 0 | 0 |
| gnalloy epoll tcp-echo 1KiB | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 5012 | 0 | 0 |
| gnalloy epoll tcp-echo 1KiB | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 4761 | 0 | 0 |
| gnalloy epoll tcp-echo 1KiB | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 4155 | 0 | 0 |
| netty epoll tcp-echo 1KiB | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 6010 | 0 | 0 |
| netty epoll tcp-echo 1KiB | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 6678 | 0 | 0 |
| netty epoll tcp-echo 1KiB | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 6422 | 0 | 0 |
| gnet poller tcp-echo 1KiB | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 6447 | 0 | 0 |
| gnet poller tcp-echo 1KiB | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 6358 | 0 | 0 |
| gnet poller tcp-echo 1KiB | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 6215 | 0 | 0 |
| netpoll poller tcp-echo 1KiB | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 4806 | 0 | 0 |
| netpoll poller tcp-echo 1KiB | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 4873 | 0 | 0 |
| netpoll poller tcp-echo 1KiB | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 4488 | 0 | 0 |
| gnalloy epoll tcp-echo 16KiB | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 6221 | 0 | 0 |
| gnalloy epoll tcp-echo 16KiB | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 7260 | 0 | 0 |
| gnalloy epoll tcp-echo 16KiB | gnalloy | tcp-echo | BenchmarkGnalloyTCPEcho-8 | 6678 | 0 | 0 |
| netty epoll tcp-echo 16KiB | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 7862 | 0 | 0 |
| netty epoll tcp-echo 16KiB | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 8681 | 0 | 0 |
| netty epoll tcp-echo 16KiB | netty | tcp-echo | BenchmarkNettyTCPEcho-8 | 11134 | 0 | 0 |
| gnet poller tcp-echo 16KiB | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 8596 | 0 | 0 |
| gnet poller tcp-echo 16KiB | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 7538 | 0 | 0 |
| gnet poller tcp-echo 16KiB | gnet | tcp-echo | BenchmarkGnetTCPEcho-8 | 8477 | 0 | 0 |
| netpoll poller tcp-echo 16KiB | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 5970 | 0 | 0 |
| netpoll poller tcp-echo 16KiB | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 6268 | 0 | 0 |
| netpoll poller tcp-echo 16KiB | netpoll | tcp-echo | BenchmarkNetpollTCPEcho-8 | 6717 | 0 | 0 |
| gnalloy epoll udp-echo 128B | gnalloy | udp-echo | BenchmarkGnalloyUDPEcho-8 | 8965 | 0 | 0 |
| gnalloy epoll udp-echo 128B | gnalloy | udp-echo | BenchmarkGnalloyUDPEcho-8 | 10006 | 0 | 0 |
| gnalloy epoll udp-echo 128B | gnalloy | udp-echo | BenchmarkGnalloyUDPEcho-8 | 8951 | 0 | 0 |
| netty epoll udp-echo 128B | netty | udp-echo | BenchmarkNettyUDPEcho-8 | 10521 | 0 | 0 |
| netty epoll udp-echo 128B | netty | udp-echo | BenchmarkNettyUDPEcho-8 | 10463 | 0 | 0 |
| netty epoll udp-echo 128B | netty | udp-echo | BenchmarkNettyUDPEcho-8 | 10432 | 0 | 0 |
| gnet poller udp-echo 128B | gnet | udp-echo | BenchmarkGnetUDPEcho-8 | 5661 | 0 | 0 |
| gnet poller udp-echo 128B | gnet | udp-echo | BenchmarkGnetUDPEcho-8 | 6047 | 0 | 0 |
| gnet poller udp-echo 128B | gnet | udp-echo | BenchmarkGnetUDPEcho-8 | 5868 | 0 | 0 |
| gnalloy epoll udp-echo 1KiB | gnalloy | udp-echo | BenchmarkGnalloyUDPEcho-8 | 10224 | 0 | 0 |
| gnalloy epoll udp-echo 1KiB | gnalloy | udp-echo | BenchmarkGnalloyUDPEcho-8 | 9367 | 0 | 0 |
| gnalloy epoll udp-echo 1KiB | gnalloy | udp-echo | BenchmarkGnalloyUDPEcho-8 | 10107 | 0 | 0 |
| netty epoll udp-echo 1KiB | netty | udp-echo | BenchmarkNettyUDPEcho-8 | 11017 | 0 | 0 |
| netty epoll udp-echo 1KiB | netty | udp-echo | BenchmarkNettyUDPEcho-8 | 10917 | 0 | 0 |
| netty epoll udp-echo 1KiB | netty | udp-echo | BenchmarkNettyUDPEcho-8 | 10850 | 0 | 0 |
| gnet poller udp-echo 1KiB | gnet | udp-echo | BenchmarkGnetUDPEcho-8 | 6477 | 0 | 0 |
| gnet poller udp-echo 1KiB | gnet | udp-echo | BenchmarkGnetUDPEcho-8 | 5942 | 0 | 0 |
| gnet poller udp-echo 1KiB | gnet | udp-echo | BenchmarkGnetUDPEcho-8 | 6026 | 0 | 0 |
| gnalloy epoll http1 128B | gnalloy | http1 | BenchmarkGnalloyHTTP1-8 | 4658 | 0 | 0 |
| gnalloy epoll http1 128B | gnalloy | http1 | BenchmarkGnalloyHTTP1-8 | 4359 | 0 | 0 |
| gnalloy epoll http1 128B | gnalloy | http1 | BenchmarkGnalloyHTTP1-8 | 4294 | 0 | 0 |
| netty epoll http1 128B | netty | http1 | BenchmarkNettyHTTP1-8 | 9309 | 0 | 0 |
| netty epoll http1 128B | netty | http1 | BenchmarkNettyHTTP1-8 | 8993 | 0 | 0 |
| netty epoll http1 128B | netty | http1 | BenchmarkNettyHTTP1-8 | 8205 | 0 | 0 |
| gnet poller http1 128B | gnet | http1 | BenchmarkGnetHTTP1-8 | 6550 | 0 | 0 |
| gnet poller http1 128B | gnet | http1 | BenchmarkGnetHTTP1-8 | 6798 | 0 | 0 |
| gnet poller http1 128B | gnet | http1 | BenchmarkGnetHTTP1-8 | 6549 | 0 | 0 |
| fasthttp http1 128B | fasthttp | http1 | BenchmarkFastHTTPHTTP1-8 | 5707 | 0 | 0 |
| fasthttp http1 128B | fasthttp | http1 | BenchmarkFastHTTPHTTP1-8 | 6570 | 0 | 0 |
| fasthttp http1 128B | fasthttp | http1 | BenchmarkFastHTTPHTTP1-8 | 5755 | 0 | 0 |
| netpoll poller http1 128B | netpoll | http1 | BenchmarkNetpollHTTP1-8 | 4394 | 0 | 0 |
| netpoll poller http1 128B | netpoll | http1 | BenchmarkNetpollHTTP1-8 | 4871 | 0 | 0 |
| netpoll poller http1 128B | netpoll | http1 | BenchmarkNetpollHTTP1-8 | 4831 | 0 | 0 |
| gnalloy epoll http1 1KiB | gnalloy | http1 | BenchmarkGnalloyHTTP1-8 | 4347 | 0 | 0 |
| gnalloy epoll http1 1KiB | gnalloy | http1 | BenchmarkGnalloyHTTP1-8 | 4342 | 0 | 0 |
| gnalloy epoll http1 1KiB | gnalloy | http1 | BenchmarkGnalloyHTTP1-8 | 5092 | 0 | 0 |
| netty epoll http1 1KiB | netty | http1 | BenchmarkNettyHTTP1-8 | 10020 | 0 | 0 |
| netty epoll http1 1KiB | netty | http1 | BenchmarkNettyHTTP1-8 | 10796 | 0 | 0 |
| netty epoll http1 1KiB | netty | http1 | BenchmarkNettyHTTP1-8 | 9222 | 0 | 0 |
| gnet poller http1 1KiB | gnet | http1 | BenchmarkGnetHTTP1-8 | 6721 | 0 | 0 |
| gnet poller http1 1KiB | gnet | http1 | BenchmarkGnetHTTP1-8 | 6677 | 0 | 0 |
| gnet poller http1 1KiB | gnet | http1 | BenchmarkGnetHTTP1-8 | 6650 | 0 | 0 |
| fasthttp http1 1KiB | fasthttp | http1 | BenchmarkFastHTTPHTTP1-8 | 5713 | 0 | 0 |
| fasthttp http1 1KiB | fasthttp | http1 | BenchmarkFastHTTPHTTP1-8 | 6747 | 0 | 0 |
| fasthttp http1 1KiB | fasthttp | http1 | BenchmarkFastHTTPHTTP1-8 | 5760 | 0 | 0 |
| netpoll poller http1 1KiB | netpoll | http1 | BenchmarkNetpollHTTP1-8 | 4898 | 0 | 0 |
| netpoll poller http1 1KiB | netpoll | http1 | BenchmarkNetpollHTTP1-8 | 4862 | 0 | 0 |
| netpoll poller http1 1KiB | netpoll | http1 | BenchmarkNetpollHTTP1-8 | 4686 | 0 | 0 |
| gnalloy epoll https1 tls11 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 12645 | 0 | 0 |
| gnalloy epoll https1 tls11 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 12511 | 0 | 0 |
| gnalloy epoll https1 tls11 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 11081 | 0 | 0 |
| netty epoll https1 tls11 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 17463 | 0 | 0 |
| netty epoll https1 tls11 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 17120 | 0 | 0 |
| netty epoll https1 tls11 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 14276 | 0 | 0 |
| fasthttp https1 tls11 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 9237 | 0 | 0 |
| fasthttp https1 tls11 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 8981 | 0 | 0 |
| fasthttp https1 tls11 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 8248 | 0 | 0 |
| gnalloy epoll https1 tls11 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 12382 | 0 | 0 |
| gnalloy epoll https1 tls11 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 14736 | 0 | 0 |
| gnalloy epoll https1 tls11 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 12367 | 0 | 0 |
| netty epoll https1 tls11 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 18045 | 0 | 0 |
| netty epoll https1 tls11 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 17614 | 0 | 0 |
| netty epoll https1 tls11 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 16841 | 0 | 0 |
| fasthttp https1 tls11 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 11520 | 0 | 0 |
| fasthttp https1 tls11 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 9404 | 0 | 0 |
| fasthttp https1 tls11 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 11150 | 0 | 0 |
| gnalloy epoll https1 tls12 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 9272 | 0 | 0 |
| gnalloy epoll https1 tls12 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 7747 | 0 | 0 |
| gnalloy epoll https1 tls12 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 7767 | 0 | 0 |
| netty epoll https1 tls12 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 15565 | 0 | 0 |
| netty epoll https1 tls12 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 15179 | 0 | 0 |
| netty epoll https1 tls12 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 14855 | 0 | 0 |
| fasthttp https1 tls12 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 7422 | 0 | 0 |
| fasthttp https1 tls12 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 6236 | 0 | 0 |
| fasthttp https1 tls12 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 5408 | 0 | 0 |
| gnalloy epoll https1 tls12 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 9449 | 0 | 0 |
| gnalloy epoll https1 tls12 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 9606 | 0 | 0 |
| gnalloy epoll https1 tls12 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 9535 | 0 | 0 |
| netty epoll https1 tls12 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 15201 | 0 | 0 |
| netty epoll https1 tls12 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 16267 | 0 | 0 |
| netty epoll https1 tls12 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 15850 | 0 | 0 |
| fasthttp https1 tls12 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 6400 | 0 | 0 |
| fasthttp https1 tls12 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 6474 | 0 | 0 |
| fasthttp https1 tls12 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 5591 | 0 | 0 |
| gnalloy epoll https1 tls13 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 9045 | 0 | 0 |
| gnalloy epoll https1 tls13 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 9235 | 0 | 0 |
| gnalloy epoll https1 tls13 128B | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 10512 | 0 | 0 |
| netty epoll https1 tls13 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 18444 | 0 | 0 |
| netty epoll https1 tls13 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 14763 | 0 | 0 |
| netty epoll https1 tls13 128B | netty | https1 | BenchmarkNettyHTTPS1-8 | 17319 | 0 | 0 |
| fasthttp https1 tls13 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 6492 | 0 | 0 |
| fasthttp https1 tls13 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 6122 | 0 | 0 |
| fasthttp https1 tls13 128B | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 5476 | 0 | 0 |
| gnalloy epoll https1 tls13 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 8658 | 0 | 0 |
| gnalloy epoll https1 tls13 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 11479 | 0 | 0 |
| gnalloy epoll https1 tls13 1KiB | gnalloy | https1 | BenchmarkGnalloyHTTPS1-8 | 9554 | 0 | 0 |
| netty epoll https1 tls13 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 18735 | 0 | 0 |
| netty epoll https1 tls13 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 20867 | 0 | 0 |
| netty epoll https1 tls13 1KiB | netty | https1 | BenchmarkNettyHTTPS1-8 | 18081 | 0 | 0 |
| fasthttp https1 tls13 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 5720 | 0 | 0 |
| fasthttp https1 tls13 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 5661 | 0 | 0 |
| fasthttp https1 tls13 1KiB | fasthttp | https1 | BenchmarkFastHTTPHTTPS1-8 | 6191 | 0 | 0 |
| gnalloy epoll http2 128B | gnalloy | http2 | BenchmarkGnalloyHTTP2-8 | 8992 | 0 | 0 |
| gnalloy epoll http2 128B | gnalloy | http2 | BenchmarkGnalloyHTTP2-8 | 10298 | 0 | 0 |
| gnalloy epoll http2 128B | gnalloy | http2 | BenchmarkGnalloyHTTP2-8 | 9611 | 0 | 0 |
| netty epoll http2 128B | netty | http2 | BenchmarkNettyHTTP2-8 | 15459 | 0 | 0 |
| netty epoll http2 128B | netty | http2 | BenchmarkNettyHTTP2-8 | 17414 | 0 | 0 |
| netty epoll http2 128B | netty | http2 | BenchmarkNettyHTTP2-8 | 13831 | 0 | 0 |
| gnalloy epoll http2 1KiB | gnalloy | http2 | BenchmarkGnalloyHTTP2-8 | 10701 | 0 | 0 |
| gnalloy epoll http2 1KiB | gnalloy | http2 | BenchmarkGnalloyHTTP2-8 | 9065 | 0 | 0 |
| gnalloy epoll http2 1KiB | gnalloy | http2 | BenchmarkGnalloyHTTP2-8 | 10336 | 0 | 0 |
| netty epoll http2 1KiB | netty | http2 | BenchmarkNettyHTTP2-8 | 15261 | 0 | 0 |
| netty epoll http2 1KiB | netty | http2 | BenchmarkNettyHTTP2-8 | 15104 | 0 | 0 |
| netty epoll http2 1KiB | netty | http2 | BenchmarkNettyHTTP2-8 | 13903 | 0 | 0 |
| gnalloy epoll https2 tls12 128B | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 16460 | 0 | 0 |
| gnalloy epoll https2 tls12 128B | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 15015 | 0 | 0 |
| gnalloy epoll https2 tls12 128B | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 16685 | 0 | 0 |
| netty epoll https2 tls12 128B | netty | https2 | BenchmarkNettyHTTPS2-8 | 20900 | 0 | 0 |
| netty epoll https2 tls12 128B | netty | https2 | BenchmarkNettyHTTPS2-8 | 23008 | 0 | 0 |
| netty epoll https2 tls12 128B | netty | https2 | BenchmarkNettyHTTPS2-8 | 23747 | 0 | 0 |
| gnalloy epoll https2 tls12 1KiB | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 17406 | 0 | 0 |
| gnalloy epoll https2 tls12 1KiB | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 15187 | 0 | 0 |
| gnalloy epoll https2 tls12 1KiB | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 17033 | 0 | 0 |
| netty epoll https2 tls12 1KiB | netty | https2 | BenchmarkNettyHTTPS2-8 | 21740 | 0 | 0 |
| netty epoll https2 tls12 1KiB | netty | https2 | BenchmarkNettyHTTPS2-8 | 22893 | 0 | 0 |
| netty epoll https2 tls12 1KiB | netty | https2 | BenchmarkNettyHTTPS2-8 | 23129 | 0 | 0 |
| gnalloy epoll https2 tls13 128B | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 16475 | 0 | 0 |
| gnalloy epoll https2 tls13 128B | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 15274 | 0 | 0 |
| gnalloy epoll https2 tls13 128B | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 16879 | 0 | 0 |
| netty epoll https2 tls13 128B | netty | https2 | BenchmarkNettyHTTPS2-8 | 25867 | 0 | 0 |
| netty epoll https2 tls13 128B | netty | https2 | BenchmarkNettyHTTPS2-8 | 24564 | 0 | 0 |
| netty epoll https2 tls13 128B | netty | https2 | BenchmarkNettyHTTPS2-8 | 25690 | 0 | 0 |
| gnalloy epoll https2 tls13 1KiB | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 16554 | 0 | 0 |
| gnalloy epoll https2 tls13 1KiB | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 17725 | 0 | 0 |
| gnalloy epoll https2 tls13 1KiB | gnalloy | https2 | BenchmarkGnalloyHTTPS2-8 | 15047 | 0 | 0 |
| netty epoll https2 tls13 1KiB | netty | https2 | BenchmarkNettyHTTPS2-8 | 23878 | 0 | 0 |
| netty epoll https2 tls13 1KiB | netty | https2 | BenchmarkNettyHTTPS2-8 | 26217 | 0 | 0 |
| netty epoll https2 tls13 1KiB | netty | https2 | BenchmarkNettyHTTPS2-8 | 27466 | 0 | 0 |
| gnalloy rfc9000 http3 tls13 128B | gnalloy | http3 | BenchmarkGnalloyHTTP3-8 | 35429 | 0 | 0 |
| gnalloy rfc9000 http3 tls13 128B | gnalloy | http3 | BenchmarkGnalloyHTTP3-8 | 33941 | 0 | 0 |
| gnalloy rfc9000 http3 tls13 128B | gnalloy | http3 | BenchmarkGnalloyHTTP3-8 | 33230 | 0 | 0 |
| netty epoll http3 tls13 128B | netty | http3 | BenchmarkNettyHTTP3-8 | 77078 | 0 | 0 |
| netty epoll http3 tls13 128B | netty | http3 | BenchmarkNettyHTTP3-8 | 76787 | 0 | 0 |
| netty epoll http3 tls13 128B | netty | http3 | BenchmarkNettyHTTP3-8 | 78318 | 0 | 0 |
| gnalloy rfc9000 http3 tls13 1KiB | gnalloy | http3 | BenchmarkGnalloyHTTP3-8 | 34842 | 0 | 0 |
| gnalloy rfc9000 http3 tls13 1KiB | gnalloy | http3 | BenchmarkGnalloyHTTP3-8 | 32437 | 0 | 0 |
| gnalloy rfc9000 http3 tls13 1KiB | gnalloy | http3 | BenchmarkGnalloyHTTP3-8 | 32473 | 0 | 0 |
| netty epoll http3 tls13 1KiB | netty | http3 | BenchmarkNettyHTTP3-8 | 77212 | 0 | 0 |
| netty epoll http3 tls13 1KiB | netty | http3 | BenchmarkNettyHTTP3-8 | 78150 | 0 | 0 |
| netty epoll http3 tls13 1KiB | netty | http3 | BenchmarkNettyHTTP3-8 | 78712 | 0 | 0 |
| gnalloy rfc9000 quic-stream tls13 128B | gnalloy | quic-stream | BenchmarkGnalloyQUICStream-8 | 23169 | 0 | 0 |
| gnalloy rfc9000 quic-stream tls13 128B | gnalloy | quic-stream | BenchmarkGnalloyQUICStream-8 | 22088 | 0 | 0 |
| gnalloy rfc9000 quic-stream tls13 128B | gnalloy | quic-stream | BenchmarkGnalloyQUICStream-8 | 21038 | 0 | 0 |
| netty epoll quic-stream tls13 128B | netty | quic-stream | BenchmarkNettyQUICStream-8 | 1776595 | 0 | 0 |
| netty epoll quic-stream tls13 128B | netty | quic-stream | BenchmarkNettyQUICStream-8 | 1781330 | 0 | 0 |
| netty epoll quic-stream tls13 128B | netty | quic-stream | BenchmarkNettyQUICStream-8 | 1787604 | 0 | 0 |
| gnalloy rfc9000 quic-stream tls13 1KiB | gnalloy | quic-stream | BenchmarkGnalloyQUICStream-8 | 21145 | 0 | 0 |
| gnalloy rfc9000 quic-stream tls13 1KiB | gnalloy | quic-stream | BenchmarkGnalloyQUICStream-8 | 21209 | 0 | 0 |
| gnalloy rfc9000 quic-stream tls13 1KiB | gnalloy | quic-stream | BenchmarkGnalloyQUICStream-8 | 22191 | 0 | 0 |
| netty epoll quic-stream tls13 1KiB | netty | quic-stream | BenchmarkNettyQUICStream-8 | 1479702 | 0 | 0 |
| netty epoll quic-stream tls13 1KiB | netty | quic-stream | BenchmarkNettyQUICStream-8 | 1505930 | 0 | 0 |
| netty epoll quic-stream tls13 1KiB | netty | quic-stream | BenchmarkNettyQUICStream-8 | 1500683 | 0 | 0 |


## Scenarios

### gnalloy epoll tcp-echo 64B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | tcp-echo |
| backend | epoll |
| payload | 64B |
| warmup | 1 |
| repeat | 3 |
| duration | 7.88923237s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol tcp-echo -backend epoll -payload 64 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=223165 p95LatencyNs=866339 p99LatencyNs=2098731 p999LatencyNs=4401147 maxLatencyNs=6406094 rssBytes=17702912 heapAllocBytes=2451680 heapSysBytes=7602176 heapObjects=4068 gcCount=0 gcPauseNs=0 goroutines=8 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=1.635026519s throughput=195715.48 ops/s
BenchmarkGnalloyTCPEcho-8 320000 5109 ns/op

framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=236056 p95LatencyNs=1240174 p99LatencyNs=2503377 p999LatencyNs=5006172 maxLatencyNs=9125494 rssBytes=17567744 heapAllocBytes=2458880 heapSysBytes=7569408 heapObjects=4050 gcCount=0 gcPauseNs=0 goroutines=8 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=2.022418385s throughput=158226.41 ops/s
BenchmarkGnalloyTCPEcho-8 320000 6320 ns/op

framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=229367 p95LatencyNs=1046689 p99LatencyNs=2213392 p999LatencyNs=4042748 maxLatencyNs=6321329 rssBytes=19808256 heapAllocBytes=2434320 heapSysBytes=11796480 heapObjects=3989 gcCount=0 gcPauseNs=0 goroutines=8 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=1.731941271s throughput=184763.77 ops/s
BenchmarkGnalloyTCPEcho-8 320000 5412 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.877799581s | 1 | 1 |  |
| 2 | 0 | 2.266393425s | 1 | 1 |  |
| 3 | 0 | 1.992202883s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 64 | 64 | 5000 | 320000 | 0 | 195715.48 | 223165 | 2098731 | 17702912 | 0 | 1.635026519s |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 64 | 64 | 5000 | 320000 | 0 | 158226.41 | 236056 | 2503377 | 17567744 | 0 | 2.022418385s |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 64 | 64 | 5000 | 320000 | 0 | 184763.77 | 229367 | 2213392 | 19808256 | 0 | 1.731941271s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 5109 | 0 | 0 |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 6320 | 0 | 0 |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 5412 | 0 | 0 |

### netty epoll tcp-echo 64B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | tcp-echo |
| backend | epoll |
| payload | 64B |
| warmup | 1 |
| repeat | 3 |
| duration | 23.750714932s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol tcp-echo --backend epoll --payload 64 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=222599 p95LatencyNs=1238452 p99LatencyNs=3253290 p999LatencyNs=7152891 maxLatencyNs=11283619 rssBytes=174456832 heapAllocBytes=48238888 heapSysBytes=264241152 heapObjects=0 gcCount=0 gcPauseNs=0 goroutines=0 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.279780891S throughput=140364.37 ops/s
BenchmarkNettyTCPEcho-8 320000 7124 ns/op

framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=232004 p95LatencyNs=958861 p99LatencyNs=2323962 p999LatencyNs=5219790 maxLatencyNs=8328897 rssBytes=184086528 heapAllocBytes=48886944 heapSysBytes=264241152 heapObjects=0 gcCount=0 gcPauseNs=0 goroutines=0 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=PT1.831358471S throughput=174733.68 ops/s
BenchmarkNettyTCPEcho-8 320000 5723 ns/op

framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=248350 p95LatencyNs=887945 p99LatencyNs=1790865 p999LatencyNs=4729860 maxLatencyNs=6746121 rssBytes=187568128 heapAllocBytes=48017160 heapSysBytes=264241152 heapObjects=0 gcCount=0 gcPauseNs=0 goroutines=0 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=PT1.73890027S throughput=184024.35 ops/s
BenchmarkNettyTCPEcho-8 320000 5434 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 6.217040813s | 1 | 1 |  |
| 2 | 0 | 5.810420148s | 1 | 1 |  |
| 3 | 0 | 5.801369731s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | tcp-echo |  | epoll | eventLoops=8 | 64 | 64 | 5000 | 320000 | 0 | 140364.37 | 222599 | 3253290 | 174456832 | 0 | 2.279780891s |
| netty | tcp-echo |  | epoll | eventLoops=8 | 64 | 64 | 5000 | 320000 | 0 | 174733.68 | 232004 | 2323962 | 184086528 | 0 | 1.831358471s |
| netty | tcp-echo |  | epoll | eventLoops=8 | 64 | 64 | 5000 | 320000 | 0 | 184024.35 | 248350 | 1790865 | 187568128 | 0 | 1.73890027s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyTCPEcho-8 | 320000 | 7124 | 0 | 0 |
| BenchmarkNettyTCPEcho-8 | 320000 | 5723 | 0 | 0 |
| BenchmarkNettyTCPEcho-8 | 320000 | 5434 | 0 | 0 |

### gnet poller tcp-echo 64B

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | tcp-echo |
| backend | poller |
| payload | 64B |
| warmup | 1 |
| repeat | 3 |
| duration | 11.351228492s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnet-bench -protocol tcp-echo -payload 64 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=269542 p95LatencyNs=1119386 p99LatencyNs=1915505 p999LatencyNs=3741108 maxLatencyNs=5941741 rssBytes=16527360 heapAllocBytes=1221960 heapSysBytes=11501568 heapObjects=3109 gcCount=0 gcPauseNs=0 goroutines=13 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=2.003436944s throughput=159725.52 ops/s
BenchmarkGnetTCPEcho-8 320000 6261 ns/op

framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=283287 p95LatencyNs=1016768 p99LatencyNs=1871223 p999LatencyNs=3114313 maxLatencyNs=3925574 rssBytes=16490496 heapAllocBytes=1216104 heapSysBytes=7405568 heapObjects=3151 gcCount=0 gcPauseNs=0 goroutines=13 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=1.960210345s throughput=163247.79 ops/s
BenchmarkGnetTCPEcho-8 320000 6126 ns/op

framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=257061 p95LatencyNs=1148376 p99LatencyNs=2007350 p999LatencyNs=3601261 maxLatencyNs=5205991 rssBytes=16637952 heapAllocBytes=1237712 heapSysBytes=7340032 heapObjects=3217 gcCount=0 gcPauseNs=0 goroutines=13 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=2.001886435s throughput=159849.23 ops/s
BenchmarkGnetTCPEcho-8 320000 6256 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.786583045s | 1 | 1 |  |
| 2 | 0 | 2.839230232s | 1 | 1 |  |
| 3 | 0 | 2.880464914s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnet | tcp-echo |  | poller | eventLoops=8 | 64 | 64 | 5000 | 320000 | 0 | 159725.52 | 269542 | 1915505 | 16527360 | 0 | 2.003436944s |
| gnet | tcp-echo |  | poller | eventLoops=8 | 64 | 64 | 5000 | 320000 | 0 | 163247.79 | 283287 | 1871223 | 16490496 | 0 | 1.960210345s |
| gnet | tcp-echo |  | poller | eventLoops=8 | 64 | 64 | 5000 | 320000 | 0 | 159849.23 | 257061 | 2007350 | 16637952 | 0 | 2.001886435s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnetTCPEcho-8 | 320000 | 6261 | 0 | 0 |
| BenchmarkGnetTCPEcho-8 | 320000 | 6126 | 0 | 0 |
| BenchmarkGnetTCPEcho-8 | 320000 | 6256 | 0 | 0 |

### netpoll poller tcp-echo 64B

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | tcp-echo |
| backend | poller |
| payload | 64B |
| warmup | 1 |
| repeat | 3 |
| duration | 6.653400934s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/netpoll-bench -protocol tcp-echo -payload 64 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=250792 p95LatencyNs=508770 p99LatencyNs=813925 p999LatencyNs=936610 maxLatencyNs=1355985 rssBytes=19472384 heapAllocBytes=4225136 heapSysBytes=10321920 heapObjects=66574 gcCount=10 gcPauseNs=2091586 goroutines=5 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=1.403922257s throughput=227932.85 ops/s
BenchmarkNetpollTCPEcho-8 320000 4387 ns/op

framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=252739 p95LatencyNs=486669 p99LatencyNs=591290 p999LatencyNs=885902 maxLatencyNs=2705035 rssBytes=19181568 heapAllocBytes=4257000 heapSysBytes=14876672 heapObjects=68774 gcCount=11 gcPauseNs=1757726 goroutines=5 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=1.405705728s throughput=227643.66 ops/s
BenchmarkNetpollTCPEcho-8 320000 4393 ns/op

framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=252931 p95LatencyNs=491796 p99LatencyNs=602507 p999LatencyNs=1023504 maxLatencyNs=1773488 rssBytes=17031168 heapAllocBytes=4044288 heapSysBytes=10584064 heapObjects=59186 gcCount=11 gcPauseNs=1868871 goroutines=3 payload=64 connections=64 messages=5000 total=320000 errors=0 elapsed=1.406036888s throughput=227590.05 ops/s
BenchmarkNetpollTCPEcho-8 320000 4394 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.621010355s | 1 | 1 |  |
| 2 | 0 | 1.689033195s | 1 | 1 |  |
| 3 | 0 | 1.633088102s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netpoll | tcp-echo |  | poller | - | 64 | 64 | 5000 | 320000 | 0 | 227932.85 | 250792 | 813925 | 19472384 | 10 | 1.403922257s |
| netpoll | tcp-echo |  | poller | - | 64 | 64 | 5000 | 320000 | 0 | 227643.66 | 252739 | 591290 | 19181568 | 11 | 1.405705728s |
| netpoll | tcp-echo |  | poller | - | 64 | 64 | 5000 | 320000 | 0 | 227590.05 | 252931 | 602507 | 17031168 | 11 | 1.406036888s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNetpollTCPEcho-8 | 320000 | 4387 | 0 | 0 |
| BenchmarkNetpollTCPEcho-8 | 320000 | 4393 | 0 | 0 |
| BenchmarkNetpollTCPEcho-8 | 320000 | 4394 | 0 | 0 |

### gnalloy epoll tcp-echo 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | tcp-echo |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 6.910634869s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol tcp-echo -backend epoll -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=196840 p95LatencyNs=866277 p99LatencyNs=2204752 p999LatencyNs=4438981 maxLatencyNs=9743226 rssBytes=17604608 heapAllocBytes=2569672 heapSysBytes=11829248 heapObjects=4033 gcCount=0 gcPauseNs=0 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.60399466s throughput=199501.91 ops/s
BenchmarkGnalloyTCPEcho-8 320000 5012 ns/op

framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=196072 p95LatencyNs=808457 p99LatencyNs=2038075 p999LatencyNs=4371609 maxLatencyNs=5819622 rssBytes=17707008 heapAllocBytes=2559896 heapSysBytes=11829248 heapObjects=4014 gcCount=0 gcPauseNs=0 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.523563855s throughput=210033.86 ops/s
BenchmarkGnalloyTCPEcho-8 320000 4761 ns/op

framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=193720 p95LatencyNs=626387 p99LatencyNs=1257934 p999LatencyNs=3321577 maxLatencyNs=4397048 rssBytes=17440768 heapAllocBytes=2559984 heapSysBytes=7634944 heapObjects=4035 gcCount=0 gcPauseNs=0 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.329514256s throughput=240689.41 ops/s
BenchmarkGnalloyTCPEcho-8 320000 4155 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.820728677s | 1 | 1 |  |
| 2 | 0 | 1.761418731s | 1 | 1 |  |
| 3 | 0 | 1.581860171s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 199501.91 | 196840 | 2204752 | 17604608 | 0 | 1.60399466s |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 210033.86 | 196072 | 2038075 | 17707008 | 0 | 1.523563855s |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 240689.41 | 193720 | 1257934 | 17440768 | 0 | 1.329514256s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 5012 | 0 | 0 |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 4761 | 0 | 0 |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 4155 | 0 | 0 |

### netty epoll tcp-echo 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | tcp-echo |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 24.29817556s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol tcp-echo --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=199385 p95LatencyNs=1072173 p99LatencyNs=2779887 p999LatencyNs=6229234 maxLatencyNs=8083663 rssBytes=190599168 heapAllocBytes=49084520 heapSysBytes=264241152 heapObjects=0 gcCount=0 gcPauseNs=0 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT1.923185676S throughput=166390.59 ops/s
BenchmarkNettyTCPEcho-8 320000 6010 ns/op

framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=202064 p95LatencyNs=1333387 p99LatencyNs=4140894 p999LatencyNs=8255938 maxLatencyNs=13140550 rssBytes=184356864 heapAllocBytes=48369600 heapSysBytes=264241152 heapObjects=0 gcCount=0 gcPauseNs=0 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.136875284S throughput=149751.37 ops/s
BenchmarkNettyTCPEcho-8 320000 6678 ns/op

framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=217772 p95LatencyNs=1086377 p99LatencyNs=2667913 p999LatencyNs=5512776 maxLatencyNs=12432039 rssBytes=183382016 heapAllocBytes=48895352 heapSysBytes=264241152 heapObjects=0 gcCount=0 gcPauseNs=0 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.055054523S throughput=155713.63 ops/s
BenchmarkNettyTCPEcho-8 320000 6422 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 5.892954941s | 1 | 1 |  |
| 2 | 0 | 6.332849075s | 1 | 1 |  |
| 3 | 0 | 6.035592148s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | tcp-echo |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 166390.59 | 199385 | 2779887 | 190599168 | 0 | 1.923185676s |
| netty | tcp-echo |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 149751.37 | 202064 | 4140894 | 184356864 | 0 | 2.136875283s |
| netty | tcp-echo |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 155713.63 | 217772 | 2667913 | 183382016 | 0 | 2.055054523s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyTCPEcho-8 | 320000 | 6010 | 0 | 0 |
| BenchmarkNettyTCPEcho-8 | 320000 | 6678 | 0 | 0 |
| BenchmarkNettyTCPEcho-8 | 320000 | 6422 | 0 | 0 |

### gnet poller tcp-echo 1KiB

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | tcp-echo |
| backend | poller |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 11.448077945s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnet-bench -protocol tcp-echo -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=272849 p95LatencyNs=1103303 p99LatencyNs=1825926 p999LatencyNs=3347366 maxLatencyNs=4621529 rssBytes=16490496 heapAllocBytes=1347368 heapSysBytes=11501568 heapObjects=3106 gcCount=0 gcPauseNs=0 goroutines=13 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.062934838s throughput=155118.81 ops/s
BenchmarkGnetTCPEcho-8 320000 6447 ns/op

framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=261830 p95LatencyNs=1185793 p99LatencyNs=1989227 p999LatencyNs=2973672 maxLatencyNs=5379461 rssBytes=16650240 heapAllocBytes=1352456 heapSysBytes=11534336 heapObjects=3184 gcCount=0 gcPauseNs=0 goroutines=13 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.034523732s throughput=157284.97 ops/s
BenchmarkGnetTCPEcho-8 320000 6358 ns/op

framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=292168 p95LatencyNs=1027403 p99LatencyNs=1775060 p999LatencyNs=3259783 maxLatencyNs=5113622 rssBytes=14426112 heapAllocBytes=1338712 heapSysBytes=7405568 heapObjects=3110 gcCount=0 gcPauseNs=0 goroutines=13 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.988662575s throughput=160912.16 ops/s
BenchmarkGnetTCPEcho-8 320000 6215 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.948329042s | 1 | 1 |  |
| 2 | 0 | 2.830264251s | 1 | 1 |  |
| 3 | 0 | 2.872258794s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnet | tcp-echo |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 155118.81 | 272849 | 1825926 | 16490496 | 0 | 2.062934838s |
| gnet | tcp-echo |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 157284.97 | 261830 | 1989227 | 16650240 | 0 | 2.034523732s |
| gnet | tcp-echo |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 160912.16 | 292168 | 1775060 | 14426112 | 0 | 1.988662575s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnetTCPEcho-8 | 320000 | 6447 | 0 | 0 |
| BenchmarkGnetTCPEcho-8 | 320000 | 6358 | 0 | 0 |
| BenchmarkGnetTCPEcho-8 | 320000 | 6215 | 0 | 0 |

### netpoll poller tcp-echo 1KiB

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | tcp-echo |
| backend | poller |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 6.928720115s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/netpoll-bench -protocol tcp-echo -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=256234 p95LatencyNs=625124 p99LatencyNs=1035964 p999LatencyNs=3194795 maxLatencyNs=5197174 rssBytes=21393408 heapAllocBytes=4165104 heapSysBytes=14516224 heapObjects=45710 gcCount=9 gcPauseNs=6463384 goroutines=4 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.537985602s throughput=208064.37 ops/s
BenchmarkNetpollTCPEcho-8 320000 4806 ns/op

framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=254217 p95LatencyNs=610669 p99LatencyNs=1333689 p999LatencyNs=3968459 maxLatencyNs=6222097 rssBytes=21106688 heapAllocBytes=5196096 heapSysBytes=14548992 heapObjects=91719 gcCount=9 gcPauseNs=4842037 goroutines=3 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.559239952s throughput=205228.19 ops/s
BenchmarkNetpollTCPEcho-8 320000 4873 ns/op

framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=252397 p95LatencyNs=501125 p99LatencyNs=697326 p999LatencyNs=2801885 maxLatencyNs=5268124 rssBytes=17567744 heapAllocBytes=4217440 heapSysBytes=10616832 heapObjects=69883 gcCount=9 gcPauseNs=1503243 goroutines=3 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.436014297s throughput=222839.01 ops/s
BenchmarkNetpollTCPEcho-8 320000 4488 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.758680648s | 1 | 1 |  |
| 2 | 0 | 1.794204475s | 1 | 1 |  |
| 3 | 0 | 1.661482659s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netpoll | tcp-echo |  | poller | - | 1024 | 64 | 5000 | 320000 | 0 | 208064.37 | 256234 | 1035964 | 21393408 | 9 | 1.537985602s |
| netpoll | tcp-echo |  | poller | - | 1024 | 64 | 5000 | 320000 | 0 | 205228.19 | 254217 | 1333689 | 21106688 | 9 | 1.559239952s |
| netpoll | tcp-echo |  | poller | - | 1024 | 64 | 5000 | 320000 | 0 | 222839.01 | 252397 | 697326 | 17567744 | 9 | 1.436014297s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNetpollTCPEcho-8 | 320000 | 4806 | 0 | 0 |
| BenchmarkNetpollTCPEcho-8 | 320000 | 4873 | 0 | 0 |
| BenchmarkNetpollTCPEcho-8 | 320000 | 4488 | 0 | 0 |

### gnalloy epoll tcp-echo 16KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | tcp-echo |
| backend | epoll |
| payload | 16KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 9.60753037s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol tcp-echo -backend epoll -payload 16384 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=16384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=298568 p95LatencyNs=966359 p99LatencyNs=1924686 p999LatencyNs=3937835 maxLatencyNs=10651931 rssBytes=20230144 heapAllocBytes=4918168 heapSysBytes=11796480 heapObjects=3360 gcCount=1 gcPauseNs=133187 goroutines=8 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=1.990854045s throughput=160735.04 ops/s
BenchmarkGnalloyTCPEcho-8 320000 6221 ns/op

framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=16384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=300642 p95LatencyNs=1280342 p99LatencyNs=2540461 p999LatencyNs=5430956 maxLatencyNs=8515460 rssBytes=20279296 heapAllocBytes=4919904 heapSysBytes=11763712 heapObjects=3309 gcCount=1 gcPauseNs=55383 goroutines=8 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=2.323337931s throughput=137732.87 ops/s
BenchmarkGnalloyTCPEcho-8 320000 7260 ns/op

framework=gnalloy protocol=tcp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=16384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=302102 p95LatencyNs=1092464 p99LatencyNs=2224978 p999LatencyNs=5495750 maxLatencyNs=7571392 rssBytes=20307968 heapAllocBytes=4954624 heapSysBytes=11763712 heapObjects=3349 gcCount=1 gcPauseNs=83126 goroutines=8 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=2.136989261s throughput=149743.38 ops/s
BenchmarkGnalloyTCPEcho-8 320000 6678 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.221348124s | 1 | 1 |  |
| 2 | 0 | 2.675470993s | 1 | 1 |  |
| 3 | 0 | 2.440568546s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=16384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 16384 | 64 | 5000 | 320000 | 0 | 160735.04 | 298568 | 1924686 | 20230144 | 1 | 1.990854045s |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=16384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 16384 | 64 | 5000 | 320000 | 0 | 137732.87 | 300642 | 2540461 | 20279296 | 1 | 2.323337931s |
| gnalloy | tcp-echo |  | epoll | boss=1 workers=4 readBuffer=16384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 16384 | 64 | 5000 | 320000 | 0 | 149743.38 | 302102 | 2224978 | 20307968 | 1 | 2.136989261s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 6221 | 0 | 0 |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 7260 | 0 | 0 |
| BenchmarkGnalloyTCPEcho-8 | 320000 | 6678 | 0 | 0 |

### netty epoll tcp-echo 16KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | tcp-echo |
| backend | epoll |
| payload | 16KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 28.492574871s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol tcp-echo --backend epoll --payload 16384 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=275872 p95LatencyNs=1423285 p99LatencyNs=3116394 p999LatencyNs=7040886 maxLatencyNs=12120744 rssBytes=191815680 heapAllocBytes=18153312 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=6000000 goroutines=0 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.515706819S throughput=127200.83 ops/s
BenchmarkNettyTCPEcho-8 320000 7862 ns/op

framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=306565 p95LatencyNs=1568263 p99LatencyNs=3554574 p999LatencyNs=8007990 maxLatencyNs=10246684 rssBytes=202924032 heapAllocBytes=15775680 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=8000000 goroutines=0 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.7780272S throughput=115189.66 ops/s
BenchmarkNettyTCPEcho-8 320000 8681 ns/op

framework=netty protocol=tcp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=305435 p95LatencyNs=2242926 p99LatencyNs=5016457 p999LatencyNs=11504963 maxLatencyNs=20584824 rssBytes=188547072 heapAllocBytes=17240960 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=6000000 goroutines=0 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.56287081S throughput=89815.21 ops/s
BenchmarkNettyTCPEcho-8 320000 11134 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 6.688462793s | 1 | 1 |  |
| 2 | 0 | 6.951802198s | 1 | 1 |  |
| 3 | 0 | 8.000176601s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | tcp-echo |  | epoll | eventLoops=8 | 16384 | 64 | 5000 | 320000 | 0 | 127200.83 | 275872 | 3116394 | 191815680 | 1 | 2.515706819s |
| netty | tcp-echo |  | epoll | eventLoops=8 | 16384 | 64 | 5000 | 320000 | 0 | 115189.66 | 306565 | 3554574 | 202924032 | 1 | 2.7780272s |
| netty | tcp-echo |  | epoll | eventLoops=8 | 16384 | 64 | 5000 | 320000 | 0 | 89815.21 | 305435 | 5016457 | 188547072 | 1 | 3.56287081s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyTCPEcho-8 | 320000 | 7862 | 0 | 0 |
| BenchmarkNettyTCPEcho-8 | 320000 | 8681 | 0 | 0 |
| BenchmarkNettyTCPEcho-8 | 320000 | 11134 | 0 | 0 |

### gnet poller tcp-echo 16KiB

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | tcp-echo |
| backend | poller |
| payload | 16KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 14.26871002s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnet-bench -protocol tcp-echo -payload 16384 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=357417 p95LatencyNs=1655818 p99LatencyNs=2806307 p999LatencyNs=4182208 maxLatencyNs=7130844 rssBytes=19206144 heapAllocBytes=3216008 heapSysBytes=11272192 heapObjects=2108 gcCount=1 gcPauseNs=39032 goroutines=13 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=2.750715315s throughput=116333.38 ops/s
BenchmarkGnetTCPEcho-8 320000 8596 ns/op

framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=334962 p95LatencyNs=1342371 p99LatencyNs=2251497 p999LatencyNs=3521062 maxLatencyNs=4896242 rssBytes=17113088 heapAllocBytes=3181808 heapSysBytes=7143424 heapObjects=2029 gcCount=1 gcPauseNs=336330 goroutines=13 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=2.412092759s throughput=132664.88 ops/s
BenchmarkGnetTCPEcho-8 320000 7538 ns/op

framework=gnet protocol=tcp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=337037 p95LatencyNs=1575910 p99LatencyNs=2612989 p999LatencyNs=4486438 maxLatencyNs=7387488 rssBytes=19177472 heapAllocBytes=3196664 heapSysBytes=11272192 heapObjects=2091 gcCount=1 gcPauseNs=170736 goroutines=13 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=2.712765829s throughput=117960.79 ops/s
BenchmarkGnetTCPEcho-8 320000 8477 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 3.697699415s | 1 | 1 |  |
| 2 | 0 | 3.348241579s | 1 | 1 |  |
| 3 | 0 | 3.664534301s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnet | tcp-echo |  | poller | eventLoops=8 | 16384 | 64 | 5000 | 320000 | 0 | 116333.38 | 357417 | 2806307 | 19206144 | 1 | 2.750715315s |
| gnet | tcp-echo |  | poller | eventLoops=8 | 16384 | 64 | 5000 | 320000 | 0 | 132664.88 | 334962 | 2251497 | 17113088 | 1 | 2.412092759s |
| gnet | tcp-echo |  | poller | eventLoops=8 | 16384 | 64 | 5000 | 320000 | 0 | 117960.79 | 337037 | 2612989 | 19177472 | 1 | 2.712765829s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnetTCPEcho-8 | 320000 | 8596 | 0 | 0 |
| BenchmarkGnetTCPEcho-8 | 320000 | 7538 | 0 | 0 |
| BenchmarkGnetTCPEcho-8 | 320000 | 8477 | 0 | 0 |

### netpoll poller tcp-echo 16KiB

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | tcp-echo |
| backend | poller |
| payload | 16KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 9.383621668s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/netpoll-bench -protocol tcp-echo -payload 16384 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=337738 p95LatencyNs=654941 p99LatencyNs=774620 p999LatencyNs=1392413 maxLatencyNs=3722193 rssBytes=26025984 heapAllocBytes=8856424 heapSysBytes=18874368 heapObjects=132681 gcCount=6 gcPauseNs=4048870 goroutines=4 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=1.910390456s throughput=167505.02 ops/s
BenchmarkNetpollTCPEcho-8 320000 5970 ns/op

framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=341283 p95LatencyNs=710126 p99LatencyNs=1135772 p999LatencyNs=2306251 maxLatencyNs=4327838 rssBytes=24215552 heapAllocBytes=8806448 heapSysBytes=18579456 heapObjects=112825 gcCount=7 gcPauseNs=1821614 goroutines=3 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=2.005639897s throughput=159550.08 ops/s
BenchmarkNetpollTCPEcho-8 320000 6268 ns/op

framework=netpoll protocol=tcp-echo backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=369527 p95LatencyNs=857035 p99LatencyNs=1494665 p999LatencyNs=3974304 maxLatencyNs=5571808 rssBytes=24371200 heapAllocBytes=10062640 heapSysBytes=14385152 heapObjects=176656 gcCount=6 gcPauseNs=1659573 goroutines=4 payload=16384 connections=64 messages=5000 total=320000 errors=0 elapsed=2.149485359s throughput=148872.84 ops/s
BenchmarkNetpollTCPEcho-8 320000 6717 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.272867668s | 1 | 1 |  |
| 2 | 0 | 2.369393917s | 1 | 1 |  |
| 3 | 0 | 2.458685232s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netpoll | tcp-echo |  | poller | - | 16384 | 64 | 5000 | 320000 | 0 | 167505.02 | 337738 | 774620 | 26025984 | 6 | 1.910390456s |
| netpoll | tcp-echo |  | poller | - | 16384 | 64 | 5000 | 320000 | 0 | 159550.08 | 341283 | 1135772 | 24215552 | 7 | 2.005639897s |
| netpoll | tcp-echo |  | poller | - | 16384 | 64 | 5000 | 320000 | 0 | 148872.84 | 369527 | 1494665 | 24371200 | 6 | 2.149485359s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNetpollTCPEcho-8 | 320000 | 5970 | 0 | 0 |
| BenchmarkNetpollTCPEcho-8 | 320000 | 6268 | 0 | 0 |
| BenchmarkNetpollTCPEcho-8 | 320000 | 6717 | 0 | 0 |

### fasthttp tcp echo unsupported

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | tcp-echo |
| backend | unsupported |
| payload | - |
| duration | 121ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
fasthttp is an HTTP server framework, not a raw TCP echo framework
```

### gnalloy epoll udp-echo 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | udp-echo |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 13.342792521s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol udp-echo -backend epoll -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=udp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=564348 p95LatencyNs=750059 p99LatencyNs=856062 p999LatencyNs=900590 maxLatencyNs=937331 rssBytes=18948096 heapAllocBytes=2434784 heapSysBytes=11829248 heapObjects=7866 gcCount=51 gcPauseNs=2460755 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.868726387s throughput=111547.76 ops/s
BenchmarkGnalloyUDPEcho-8 320000 8965 ns/op

framework=gnalloy protocol=udp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=623183 p95LatencyNs=934189 p99LatencyNs=2869700 p999LatencyNs=3253570 maxLatencyNs=3432419 rssBytes=18640896 heapAllocBytes=3290144 heapSysBytes=11829248 heapObjects=21753 gcCount=50 gcPauseNs=7199315 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=3.201968543s throughput=99938.52 ops/s
BenchmarkGnalloyUDPEcho-8 320000 10006 ns/op

framework=gnalloy protocol=udp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=567006 p95LatencyNs=744076 p99LatencyNs=834035 p999LatencyNs=949586 maxLatencyNs=1001711 rssBytes=19148800 heapAllocBytes=3542000 heapSysBytes=11829248 heapObjects=25922 gcCount=50 gcPauseNs=2521183 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.864341539s throughput=111718.52 ops/s
BenchmarkGnalloyUDPEcho-8 320000 8951 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 3.269406058s | 1 | 1 |  |
| 2 | 0 | 3.552858961s | 1 | 1 |  |
| 3 | 0 | 3.188355641s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 111547.76 | 564348 | 856062 | 18948096 | 51 | 2.868726387s |
| gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 99938.52 | 623183 | 2869700 | 18640896 | 50 | 3.201968543s |
| gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 111718.52 | 567006 | 834035 | 19148800 | 50 | 2.864341539s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyUDPEcho-8 | 320000 | 8965 | 0 | 0 |
| BenchmarkGnalloyUDPEcho-8 | 320000 | 10006 | 0 | 0 |
| BenchmarkGnalloyUDPEcho-8 | 320000 | 8951 | 0 | 0 |

### netty epoll udp-echo 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | udp-echo |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 29.217028093s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol udp-echo --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=udp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=656771 p95LatencyNs=820744 p99LatencyNs=903855 p999LatencyNs=1364420 maxLatencyNs=2913218 rssBytes=184332288 heapAllocBytes=93961368 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=4000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.36657383S throughput=95052.13 ops/s
BenchmarkNettyUDPEcho-8 320000 10521 ns/op

framework=netty protocol=udp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=657856 p95LatencyNs=806599 p99LatencyNs=877503 p999LatencyNs=1229331 maxLatencyNs=4502586 rssBytes=185024512 heapAllocBytes=93146824 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=4000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.348191954S throughput=95573.97 ops/s
BenchmarkNettyUDPEcho-8 320000 10463 ns/op

framework=netty protocol=udp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=647021 p95LatencyNs=807341 p99LatencyNs=900319 p999LatencyNs=1162616 maxLatencyNs=4880220 rssBytes=191569920 heapAllocBytes=95271096 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=4000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.338330821S throughput=95856.29 ops/s
BenchmarkNettyUDPEcho-8 320000 10432 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 7.527836627s | 1 | 1 |  |
| 2 | 0 | 7.509151879s | 1 | 1 |  |
| 3 | 0 | 7.58494396s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | udp-echo |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 95052.13 | 656771 | 903855 | 184332288 | 1 | 3.36657383s |
| netty | udp-echo |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 95573.97 | 657856 | 877503 | 185024512 | 1 | 3.348191954s |
| netty | udp-echo |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 95856.29 | 647021 | 900319 | 191569920 | 1 | 3.338330821s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyUDPEcho-8 | 320000 | 10521 | 0 | 0 |
| BenchmarkNettyUDPEcho-8 | 320000 | 10463 | 0 | 0 |
| BenchmarkNettyUDPEcho-8 | 320000 | 10432 | 0 | 0 |

### gnet poller udp-echo 128B

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | udp-echo |
| backend | poller |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 10.522486758s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnet-bench -protocol udp-echo -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnet protocol=udp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=277873 p95LatencyNs=921563 p99LatencyNs=1538254 p999LatencyNs=2918401 maxLatencyNs=3688531 rssBytes=17375232 heapAllocBytes=1619328 heapSysBytes=11337728 heapObjects=7913 gcCount=58 gcPauseNs=5523687 goroutines=12 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.811398831s throughput=176659.05 ops/s
BenchmarkGnetUDPEcho-8 320000 5661 ns/op

framework=gnet protocol=udp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=246300 p95LatencyNs=1071985 p99LatencyNs=1998772 p999LatencyNs=3412652 maxLatencyNs=7206565 rssBytes=18132992 heapAllocBytes=1892280 heapSysBytes=11272192 heapObjects=10783 gcCount=58 gcPauseNs=13672791 goroutines=12 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.934996819s throughput=165374.95 ops/s
BenchmarkGnetUDPEcho-8 320000 6047 ns/op

framework=gnet protocol=udp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=292895 p95LatencyNs=922422 p99LatencyNs=1520535 p999LatencyNs=2697793 maxLatencyNs=3692806 rssBytes=18169856 heapAllocBytes=1596616 heapSysBytes=11239424 heapObjects=7593 gcCount=58 gcPauseNs=5828256 goroutines=12 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.877673322s throughput=170423.68 ops/s
BenchmarkGnetUDPEcho-8 320000 5868 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.651282732s | 1 | 1 |  |
| 2 | 0 | 2.677736401s | 1 | 1 |  |
| 3 | 0 | 2.627798825s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnet | udp-echo |  | poller | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 176659.05 | 277873 | 1538254 | 17375232 | 58 | 1.811398831s |
| gnet | udp-echo |  | poller | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 165374.95 | 246300 | 1998772 | 18132992 | 58 | 1.934996819s |
| gnet | udp-echo |  | poller | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 170423.68 | 292895 | 1520535 | 18169856 | 58 | 1.877673322s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnetUDPEcho-8 | 320000 | 5661 | 0 | 0 |
| BenchmarkGnetUDPEcho-8 | 320000 | 6047 | 0 | 0 |
| BenchmarkGnetUDPEcho-8 | 320000 | 5868 | 0 | 0 |

### gnalloy epoll udp-echo 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | udp-echo |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 13.950424897s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol udp-echo -backend epoll -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=udp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=631295 p95LatencyNs=905850 p99LatencyNs=1038865 p999LatencyNs=1442302 maxLatencyNs=1841847 rssBytes=19320832 heapAllocBytes=2806184 heapSysBytes=11796480 heapObjects=11791 gcCount=48 gcPauseNs=2428473 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.271631125s throughput=97810.54 ops/s
BenchmarkGnalloyUDPEcho-8 320000 10224 ns/op

framework=gnalloy protocol=udp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=591815 p95LatencyNs=796448 p99LatencyNs=919077 p999LatencyNs=999527 maxLatencyNs=1079567 rssBytes=21143552 heapAllocBytes=3717536 heapSysBytes=11862016 heapObjects=27122 gcCount=48 gcPauseNs=2224282 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.997313711s throughput=106762.26 ops/s
BenchmarkGnalloyUDPEcho-8 320000 9367 ns/op

framework=gnalloy protocol=udp-echo backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=636763 p95LatencyNs=858640 p99LatencyNs=946924 p999LatencyNs=3102872 maxLatencyNs=3338251 rssBytes=18518016 heapAllocBytes=3380200 heapSysBytes=11862016 heapObjects=21480 gcCount=48 gcPauseNs=6148117 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.234183855s throughput=98943.05 ops/s
BenchmarkGnalloyUDPEcho-8 320000 10107 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 3.644455117s | 1 | 1 |  |
| 2 | 0 | 3.311547541s | 1 | 1 |  |
| 3 | 0 | 3.556742684s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 97810.54 | 631295 | 1038865 | 19320832 | 48 | 3.271631125s |
| gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 106762.26 | 591815 | 919077 | 21143552 | 48 | 2.997313711s |
| gnalloy | udp-echo |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 98943.05 | 636763 | 946924 | 18518016 | 48 | 3.234183855s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyUDPEcho-8 | 320000 | 10224 | 0 | 0 |
| BenchmarkGnalloyUDPEcho-8 | 320000 | 9367 | 0 | 0 |
| BenchmarkGnalloyUDPEcho-8 | 320000 | 10107 | 0 | 0 |

### netty epoll udp-echo 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | udp-echo |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 30.989911677s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol udp-echo --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=udp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=691259 p95LatencyNs=857528 p99LatencyNs=1475117 p999LatencyNs=4314971 maxLatencyNs=4776187 rssBytes=190017536 heapAllocBytes=95627352 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=4000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.525411361S throughput=90769.55 ops/s
BenchmarkNettyUDPEcho-8 320000 11017 ns/op

framework=netty protocol=udp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=679258 p95LatencyNs=844120 p99LatencyNs=986485 p999LatencyNs=3662460 maxLatencyNs=5502795 rssBytes=190365696 heapAllocBytes=95511136 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=5000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.493456251S throughput=91599.83 ops/s
BenchmarkNettyUDPEcho-8 320000 10917 ns/op

framework=netty protocol=udp-echo backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=678154 p95LatencyNs=853373 p99LatencyNs=963504 p999LatencyNs=1449409 maxLatencyNs=5669121 rssBytes=195325952 heapAllocBytes=97785600 heapSysBytes=264241152 heapObjects=0 gcCount=1 gcPauseNs=5000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.472157208S throughput=92161.73 ops/s
BenchmarkNettyUDPEcho-8 320000 10850 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 7.809982729s | 1 | 1 |  |
| 2 | 0 | 7.710780968s | 1 | 1 |  |
| 3 | 0 | 7.720860538s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | udp-echo |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 90769.55 | 691259 | 1475117 | 190017536 | 1 | 3.525411361s |
| netty | udp-echo |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 91599.83 | 679258 | 986485 | 190365696 | 1 | 3.493456251s |
| netty | udp-echo |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 92161.73 | 678154 | 963504 | 195325952 | 1 | 3.472157208s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyUDPEcho-8 | 320000 | 11017 | 0 | 0 |
| BenchmarkNettyUDPEcho-8 | 320000 | 10917 | 0 | 0 |
| BenchmarkNettyUDPEcho-8 | 320000 | 10850 | 0 | 0 |

### gnet poller udp-echo 1KiB

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | udp-echo |
| backend | poller |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 11.200964455s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnet-bench -protocol udp-echo -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnet protocol=udp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=271489 p95LatencyNs=1155388 p99LatencyNs=2133582 p999LatencyNs=4122318 maxLatencyNs=6901610 rssBytes=18018304 heapAllocBytes=2050864 heapSysBytes=11173888 heapObjects=11088 gcCount=61 gcPauseNs=10405146 goroutines=12 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.072601933s throughput=154395.30 ops/s
BenchmarkGnetUDPEcho-8 320000 6477 ns/op

framework=gnet protocol=udp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=276154 p95LatencyNs=981152 p99LatencyNs=1553677 p999LatencyNs=2556830 maxLatencyNs=3193772 rssBytes=18423808 heapAllocBytes=2499360 heapSysBytes=11337728 heapObjects=16157 gcCount=60 gcPauseNs=5126653 goroutines=12 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.901283518s throughput=168307.35 ops/s
BenchmarkGnetUDPEcho-8 320000 5942 ns/op

framework=gnet protocol=udp-echo backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=276078 p95LatencyNs=967281 p99LatencyNs=1465318 p999LatencyNs=3237328 maxLatencyNs=4863414 rssBytes=18612224 heapAllocBytes=3011224 heapSysBytes=11239424 heapObjects=21593 gcCount=60 gcPauseNs=5312952 goroutines=12 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.928183508s throughput=165959.31 ops/s
BenchmarkGnetUDPEcho-8 320000 6026 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.858710379s | 1 | 1 |  |
| 2 | 0 | 2.634990581s | 1 | 1 |  |
| 3 | 0 | 2.765170856s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnet | udp-echo |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 154395.3 | 271489 | 2133582 | 18018304 | 61 | 2.072601933s |
| gnet | udp-echo |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 168307.35 | 276154 | 1553677 | 18423808 | 60 | 1.901283518s |
| gnet | udp-echo |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 165959.31 | 276078 | 1465318 | 18612224 | 60 | 1.928183508s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnetUDPEcho-8 | 320000 | 6477 | 0 | 0 |
| BenchmarkGnetUDPEcho-8 | 320000 | 5942 | 0 | 0 |
| BenchmarkGnetUDPEcho-8 | 320000 | 6026 | 0 | 0 |

### netpoll udp echo unsupported

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | udp-echo |
| backend | unsupported |
| payload | - |
| duration | 107ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
CloudWeGo netpoll does not expose an equivalent UDP server API in this harness
```

### fasthttp udp echo unsupported

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | udp-echo |
| backend | unsupported |
| payload | - |
| duration | 83ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
fasthttp is HTTP-only and cannot execute UDP datagram echo
```

### gnalloy epoll http1 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | http1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 6.914583232s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol http1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=http1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=201750 p95LatencyNs=783744 p99LatencyNs=1816515 p999LatencyNs=4072125 maxLatencyNs=6347071 rssBytes=22425600 heapAllocBytes=3432056 heapSysBytes=15826944 heapObjects=6937 gcCount=2 gcPauseNs=267525 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.490465671s throughput=214698.00 ops/s
BenchmarkGnalloyHTTP1-8 320000 4658 ns/op

framework=gnalloy protocol=http1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=198004 p95LatencyNs=651108 p99LatencyNs=1519743 p999LatencyNs=3444861 maxLatencyNs=5084658 rssBytes=22609920 heapAllocBytes=3317408 heapSysBytes=15892480 heapObjects=3188 gcCount=2 gcPauseNs=264506 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.395037556s throughput=229384.51 ops/s
BenchmarkGnalloyHTTP1-8 320000 4359 ns/op

framework=gnalloy protocol=http1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=196303 p95LatencyNs=636589 p99LatencyNs=1340394 p999LatencyNs=4063908 maxLatencyNs=5553939 rssBytes=22482944 heapAllocBytes=3575616 heapSysBytes=11730944 heapObjects=16768 gcCount=2 gcPauseNs=152586 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.373998814s throughput=232896.85 ops/s
BenchmarkGnalloyHTTP1-8 320000 4294 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.723969679s | 1 | 1 |  |
| 2 | 0 | 1.677816588s | 1 | 1 |  |
| 3 | 0 | 1.631911117s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 214698 | 201750 | 1816515 | 22425600 | 2 | 1.490465671s |
| gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 229384.51 | 198004 | 1519743 | 22609920 | 2 | 1.395037556s |
| gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 232896.85 | 196303 | 1340394 | 22482944 | 2 | 1.373998814s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTP1-8 | 320000 | 4658 | 0 | 0 |
| BenchmarkGnalloyHTTP1-8 | 320000 | 4359 | 0 | 0 |
| BenchmarkGnalloyHTTP1-8 | 320000 | 4294 | 0 | 0 |

### netty epoll http1 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | http1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 29.062876227s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol http1 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=http1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=326256 p95LatencyNs=1675683 p99LatencyNs=3945852 p999LatencyNs=9335481 maxLatencyNs=13696535 rssBytes=372654080 heapAllocBytes=30229624 heapSysBytes=318767104 heapObjects=0 gcCount=7 gcPauseNs=39000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.978897936S throughput=107422.28 ops/s
BenchmarkNettyHTTP1-8 320000 9309 ns/op

framework=netty protocol=http1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=361435 p95LatencyNs=1502782 p99LatencyNs=3344587 p999LatencyNs=7654226 maxLatencyNs=12669138 rssBytes=418213888 heapAllocBytes=179229904 heapSysBytes=318767104 heapObjects=0 gcCount=6 gcPauseNs=26000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.877712045S throughput=111199.45 ops/s
BenchmarkNettyHTTP1-8 320000 8993 ns/op

framework=netty protocol=http1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=343733 p95LatencyNs=1446671 p99LatencyNs=3582067 p999LatencyNs=6847078 maxLatencyNs=10105297 rssBytes=322428928 heapAllocBytes=55225344 heapSysBytes=264241152 heapObjects=0 gcCount=7 gcPauseNs=17000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.625684263S throughput=121872.99 ops/s
BenchmarkNettyHTTP1-8 320000 8205 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 7.351254607s | 1 | 1 |  |
| 2 | 0 | 7.210757917s | 1 | 1 |  |
| 3 | 0 | 7.020758947s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | http1 |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 107422.28 | 326256 | 3945852 | 372654080 | 7 | 2.978897936s |
| netty | http1 |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 111199.45 | 361435 | 3344587 | 418213888 | 6 | 2.877712045s |
| netty | http1 |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 121872.99 | 343733 | 3582067 | 322428928 | 7 | 2.625684263s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTP1-8 | 320000 | 9309 | 0 | 0 |
| BenchmarkNettyHTTP1-8 | 320000 | 8993 | 0 | 0 |
| BenchmarkNettyHTTP1-8 | 320000 | 8205 | 0 | 0 |

### gnet poller http1 128B

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | http1 |
| backend | poller |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 11.858810037s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnet-bench -protocol http1 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnet protocol=http1 backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=303283 p95LatencyNs=1110460 p99LatencyNs=1832194 p999LatencyNs=2903608 maxLatencyNs=4525057 rssBytes=21409792 heapAllocBytes=4203112 heapSysBytes=15335424 heapObjects=133658 gcCount=2 gcPauseNs=133800 goroutines=13 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.095853444s throughput=152682.43 ops/s
BenchmarkGnetHTTP1-8 320000 6550 ns/op

framework=gnet protocol=http1 backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=268126 p95LatencyNs=1251367 p99LatencyNs=2373167 p999LatencyNs=5025939 maxLatencyNs=7444511 rssBytes=21635072 heapAllocBytes=2189336 heapSysBytes=11173888 heapObjects=4717 gcCount=2 gcPauseNs=249929 goroutines=13 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.17524356s throughput=147109.96 ops/s
BenchmarkGnetHTTP1-8 320000 6798 ns/op

framework=gnet protocol=http1 backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=308288 p95LatencyNs=1074326 p99LatencyNs=1545825 p999LatencyNs=2843521 maxLatencyNs=5437203 rssBytes=17272832 heapAllocBytes=4222032 heapSysBytes=11272192 heapObjects=132349 gcCount=1 gcPauseNs=47577 goroutines=13 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.095597434s throughput=152701.08 ops/s
BenchmarkGnetHTTP1-8 320000 6549 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.99606883s | 1 | 1 |  |
| 2 | 0 | 2.954500011s | 1 | 1 |  |
| 3 | 0 | 2.931432644s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnet | http1 |  | poller | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 152682.43 | 303283 | 1832194 | 21409792 | 2 | 2.095853444s |
| gnet | http1 |  | poller | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 147109.96 | 268126 | 2373167 | 21635072 | 2 | 2.17524356s |
| gnet | http1 |  | poller | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 152701.08 | 308288 | 1545825 | 17272832 | 1 | 2.095597434s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnetHTTP1-8 | 320000 | 6550 | 0 | 0 |
| BenchmarkGnetHTTP1-8 | 320000 | 6798 | 0 | 0 |
| BenchmarkGnetHTTP1-8 | 320000 | 6549 | 0 | 0 |

### fasthttp http1 128B

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | http1 |
| backend | net |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 8.550133895s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol http1 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=fasthttp protocol=http1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=294198 p95LatencyNs=749597 p99LatencyNs=1947699 p999LatencyNs=4255977 maxLatencyNs=8255863 rssBytes=19243008 heapAllocBytes=2280808 heapSysBytes=11042816 heapObjects=3592 gcCount=2 gcPauseNs=412001 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.826354799s throughput=175212.40 ops/s
BenchmarkFastHTTPHTTP1-8 320000 5707 ns/op

framework=fasthttp protocol=http1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=309032 p95LatencyNs=868203 p99LatencyNs=2605204 p999LatencyNs=5923656 maxLatencyNs=11529249 rssBytes=19161088 heapAllocBytes=4646160 heapSysBytes=10944512 heapObjects=147630 gcCount=1 gcPauseNs=3705072 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.102358265s throughput=152210.02 ops/s
BenchmarkFastHTTPHTTP1-8 320000 6570 ns/op

framework=fasthttp protocol=http1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=306408 p95LatencyNs=722794 p99LatencyNs=1693118 p999LatencyNs=4641113 maxLatencyNs=8951030 rssBytes=19116032 heapAllocBytes=4623800 heapSysBytes=11010048 heapObjects=146729 gcCount=1 gcPauseNs=137578 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.841713693s throughput=173751.22 ops/s
BenchmarkFastHTTPHTTP1-8 320000 5755 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.136342167s | 1 | 1 |  |
| 2 | 0 | 2.387960639s | 1 | 1 |  |
| 3 | 0 | 2.130077829s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | http1 |  | net | - | 128 | 64 | 5000 | 320000 | 0 | 175212.4 | 294198 | 1947699 | 19243008 | 2 | 1.826354799s |
| fasthttp | http1 |  | net | - | 128 | 64 | 5000 | 320000 | 0 | 152210.02 | 309032 | 2605204 | 19161088 | 1 | 2.102358265s |
| fasthttp | http1 |  | net | - | 128 | 64 | 5000 | 320000 | 0 | 173751.22 | 306408 | 1693118 | 19116032 | 1 | 1.841713693s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTP1-8 | 320000 | 5707 | 0 | 0 |
| BenchmarkFastHTTPHTTP1-8 | 320000 | 6570 | 0 | 0 |
| BenchmarkFastHTTPHTTP1-8 | 320000 | 5755 | 0 | 0 |

### netpoll poller http1 128B

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | http1 |
| backend | poller |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 7.114052727s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/netpoll-bench -protocol http1 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=netpoll protocol=http1 backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=244427 p95LatencyNs=485232 p99LatencyNs=656706 p999LatencyNs=1438931 maxLatencyNs=5653736 rssBytes=22401024 heapAllocBytes=3780504 heapSysBytes=18710528 heapObjects=11723 gcCount=15 gcPauseNs=4940862 goroutines=5 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.406037282s throughput=227589.98 ops/s
BenchmarkNetpollHTTP1-8 320000 4394 ns/op

framework=netpoll protocol=http1 backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=259502 p95LatencyNs=632203 p99LatencyNs=1064011 p999LatencyNs=3276265 maxLatencyNs=8109395 rssBytes=21094400 heapAllocBytes=4485648 heapSysBytes=14417920 heapObjects=27239 gcCount=14 gcPauseNs=2864775 goroutines=6 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.558716887s throughput=205297.06 ops/s
BenchmarkNetpollHTTP1-8 320000 4871 ns/op

framework=netpoll protocol=http1 backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=254878 p95LatencyNs=621560 p99LatencyNs=994058 p999LatencyNs=2770079 maxLatencyNs=4154145 rssBytes=21434368 heapAllocBytes=6770472 heapSysBytes=14254080 heapObjects=96127 gcCount=14 gcPauseNs=4071136 goroutines=3 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.545923263s throughput=206996.04 ops/s
BenchmarkNetpollHTTP1-8 320000 4831 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.701738614s | 1 | 1 |  |
| 2 | 0 | 1.860250431s | 1 | 1 |  |
| 3 | 0 | 1.81949736s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netpoll | http1 |  | poller | - | 128 | 64 | 5000 | 320000 | 0 | 227589.98 | 244427 | 656706 | 22401024 | 15 | 1.406037282s |
| netpoll | http1 |  | poller | - | 128 | 64 | 5000 | 320000 | 0 | 205297.06 | 259502 | 1064011 | 21094400 | 14 | 1.558716887s |
| netpoll | http1 |  | poller | - | 128 | 64 | 5000 | 320000 | 0 | 206996.04 | 254878 | 994058 | 21434368 | 14 | 1.545923263s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNetpollHTTP1-8 | 320000 | 4394 | 0 | 0 |
| BenchmarkNetpollHTTP1-8 | 320000 | 4871 | 0 | 0 |
| BenchmarkNetpollHTTP1-8 | 320000 | 4831 | 0 | 0 |

### gnalloy epoll http1 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | http1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 6.841408859s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol http1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=http1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=198655 p95LatencyNs=679208 p99LatencyNs=1587214 p999LatencyNs=5809583 maxLatencyNs=6077048 rssBytes=22286336 heapAllocBytes=4157128 heapSysBytes=15892480 heapObjects=42038 gcCount=2 gcPauseNs=261551 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.391095188s throughput=230034.58 ops/s
BenchmarkGnalloyHTTP1-8 320000 4347 ns/op

framework=gnalloy protocol=http1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=195524 p95LatencyNs=668443 p99LatencyNs=1368974 p999LatencyNs=3227704 maxLatencyNs=3805474 rssBytes=20144128 heapAllocBytes=4118688 heapSysBytes=11730944 heapObjects=40150 gcCount=2 gcPauseNs=362856 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.389477563s throughput=230302.39 ops/s
BenchmarkGnalloyHTTP1-8 320000 4342 ns/op

framework=gnalloy protocol=http1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=200386 p95LatencyNs=947586 p99LatencyNs=2227895 p999LatencyNs=5254476 maxLatencyNs=9081063 rssBytes=20353024 heapAllocBytes=3787496 heapSysBytes=11665408 heapObjects=17710 gcCount=2 gcPauseNs=243622 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.629516736s throughput=196377.24 ops/s
BenchmarkGnalloyHTTP1-8 320000 5092 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.669025046s | 1 | 1 |  |
| 2 | 0 | 1.679256813s | 1 | 1 |  |
| 3 | 0 | 1.865659637s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 230034.58 | 198655 | 1587214 | 22286336 | 2 | 1.391095188s |
| gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 230302.39 | 195524 | 1368974 | 20144128 | 2 | 1.389477563s |
| gnalloy | http1 |  | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 196377.24 | 200386 | 2227895 | 20353024 | 2 | 1.629516736s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTP1-8 | 320000 | 4347 | 0 | 0 |
| BenchmarkGnalloyHTTP1-8 | 320000 | 4342 | 0 | 0 |
| BenchmarkGnalloyHTTP1-8 | 320000 | 5092 | 0 | 0 |

### netty epoll http1 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | http1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 30.646391785s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol http1 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=http1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=257049 p95LatencyNs=1903629 p99LatencyNs=4267286 p999LatencyNs=11592685 maxLatencyNs=23962835 rssBytes=325349376 heapAllocBytes=112253248 heapSysBytes=264241152 heapObjects=0 gcCount=7 gcPauseNs=27000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.206305229S throughput=99803.35 ops/s
BenchmarkNettyHTTP1-8 320000 10020 ns/op

framework=netty protocol=http1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=310041 p95LatencyNs=2243040 p99LatencyNs=4964408 p999LatencyNs=8346683 maxLatencyNs=10374502 rssBytes=396431360 heapAllocBytes=108336320 heapSysBytes=318767104 heapObjects=0 gcCount=7 gcPauseNs=36000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT3.454794901S throughput=92624.89 ops/s
BenchmarkNettyHTTP1-8 320000 10796 ns/op

framework=netty protocol=http1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=326987 p95LatencyNs=1657094 p99LatencyNs=3938645 p999LatencyNs=9105131 maxLatencyNs=12155943 rssBytes=379572224 heapAllocBytes=190540384 heapSysBytes=318767104 heapObjects=0 gcCount=6 gcPauseNs=30000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT2.950980624S throughput=108438.53 ops/s
BenchmarkNettyHTTP1-8 320000 9222 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 7.82000858s | 1 | 1 |  |
| 2 | 0 | 7.806716624s | 1 | 1 |  |
| 3 | 0 | 7.36998593s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | http1 |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 99803.35 | 257049 | 4267286 | 325349376 | 7 | 3.206305229s |
| netty | http1 |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 92624.89 | 310041 | 4964408 | 396431360 | 7 | 3.454794901s |
| netty | http1 |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 108438.53 | 326987 | 3938645 | 379572224 | 6 | 2.950980624s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTP1-8 | 320000 | 10020 | 0 | 0 |
| BenchmarkNettyHTTP1-8 | 320000 | 10796 | 0 | 0 |
| BenchmarkNettyHTTP1-8 | 320000 | 9222 | 0 | 0 |

### gnet poller http1 1KiB

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | http1 |
| backend | poller |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 11.940293234s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnet-bench -protocol http1 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnet protocol=http1 backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=291419 p95LatencyNs=1149915 p99LatencyNs=1823550 p999LatencyNs=3662901 maxLatencyNs=4369390 rssBytes=21479424 heapAllocBytes=2272176 heapSysBytes=15368192 heapObjects=2673 gcCount=2 gcPauseNs=189456 goroutines=13 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.150632975s throughput=148793.40 ops/s
BenchmarkGnetHTTP1-8 320000 6721 ns/op

framework=gnet protocol=http1 backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=281141 p95LatencyNs=1262121 p99LatencyNs=2135331 p999LatencyNs=3456175 maxLatencyNs=5303720 rssBytes=21417984 heapAllocBytes=4444760 heapSysBytes=11239424 heapObjects=138795 gcCount=1 gcPauseNs=91618 goroutines=13 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.13663269s throughput=149768.37 ops/s
BenchmarkGnetHTTP1-8 320000 6677 ns/op

framework=gnet protocol=http1 backend=poller eventLoops=8 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=313952 p95LatencyNs=1115392 p99LatencyNs=1659810 p999LatencyNs=2808180 maxLatencyNs=4453391 rssBytes=19476480 heapAllocBytes=2216288 heapSysBytes=11239424 heapObjects=2191 gcCount=2 gcPauseNs=229936 goroutines=13 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.127878541s throughput=150384.52 ops/s
BenchmarkGnetHTTP1-8 320000 6650 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 3.032309536s | 1 | 1 |  |
| 2 | 0 | 2.952200387s | 1 | 1 |  |
| 3 | 0 | 3.015495194s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnet | http1 |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 148793.4 | 291419 | 1823550 | 21479424 | 2 | 2.150632975s |
| gnet | http1 |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 149768.37 | 281141 | 2135331 | 21417984 | 1 | 2.13663269s |
| gnet | http1 |  | poller | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 150384.52 | 313952 | 1659810 | 19476480 | 2 | 2.127878541s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnetHTTP1-8 | 320000 | 6721 | 0 | 0 |
| BenchmarkGnetHTTP1-8 | 320000 | 6677 | 0 | 0 |
| BenchmarkGnetHTTP1-8 | 320000 | 6650 | 0 | 0 |

### fasthttp http1 1KiB

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | http1 |
| backend | net |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 8.793629723s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol http1 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=fasthttp protocol=http1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=299096 p95LatencyNs=730508 p99LatencyNs=1442534 p999LatencyNs=4049856 maxLatencyNs=7901468 rssBytes=21114880 heapAllocBytes=4714496 heapSysBytes=15433728 heapObjects=148005 gcCount=1 gcPauseNs=273000 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.828139208s throughput=175041.37 ops/s
BenchmarkFastHTTPHTTP1-8 320000 5713 ns/op

framework=fasthttp protocol=http1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=312616 p95LatencyNs=955379 p99LatencyNs=2829956 p999LatencyNs=6424143 maxLatencyNs=10102646 rssBytes=19140608 heapAllocBytes=4757608 heapSysBytes=10977280 heapObjects=147834 gcCount=1 gcPauseNs=415173 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.158948048s throughput=148220.33 ops/s
BenchmarkFastHTTPHTTP1-8 320000 6747 ns/op

framework=fasthttp protocol=http1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=299111 p95LatencyNs=739667 p99LatencyNs=1932874 p999LatencyNs=4209230 maxLatencyNs=6517942 rssBytes=19124224 heapAllocBytes=2364368 heapSysBytes=11010048 heapObjects=3622 gcCount=2 gcPauseNs=246917 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.843120752s throughput=173618.58 ops/s
BenchmarkFastHTTPHTTP1-8 320000 5760 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.07897717s | 1 | 1 |  |
| 2 | 0 | 2.511184464s | 1 | 1 |  |
| 3 | 0 | 2.115764728s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | http1 |  | net | - | 1024 | 64 | 5000 | 320000 | 0 | 175041.37 | 299096 | 1442534 | 21114880 | 1 | 1.828139208s |
| fasthttp | http1 |  | net | - | 1024 | 64 | 5000 | 320000 | 0 | 148220.33 | 312616 | 2829956 | 19140608 | 1 | 2.158948048s |
| fasthttp | http1 |  | net | - | 1024 | 64 | 5000 | 320000 | 0 | 173618.58 | 299111 | 1932874 | 19124224 | 2 | 1.843120752s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTP1-8 | 320000 | 5713 | 0 | 0 |
| BenchmarkFastHTTPHTTP1-8 | 320000 | 6747 | 0 | 0 |
| BenchmarkFastHTTPHTTP1-8 | 320000 | 5760 | 0 | 0 |

### netpoll poller http1 1KiB

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | http1 |
| backend | poller |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 7.103471956s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/netpoll-bench -protocol http1 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=netpoll protocol=http1 backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=256833 p95LatencyNs=650877 p99LatencyNs=1177092 p999LatencyNs=2759157 maxLatencyNs=4611562 rssBytes=24276992 heapAllocBytes=7087144 heapSysBytes=18481152 heapObjects=103989 gcCount=13 gcPauseNs=3372295 goroutines=3 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.567324571s throughput=204169.58 ops/s
BenchmarkNetpollHTTP1-8 320000 4898 ns/op

framework=netpoll protocol=http1 backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=259186 p95LatencyNs=649645 p99LatencyNs=1038564 p999LatencyNs=2769124 maxLatencyNs=4244457 rssBytes=23552000 heapAllocBytes=5581552 heapSysBytes=14319616 heapObjects=52888 gcCount=13 gcPauseNs=3758756 goroutines=4 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.555948695s throughput=205662.31 ops/s
BenchmarkNetpollHTTP1-8 320000 4862 ns/op

framework=netpoll protocol=http1 backend=poller latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=252640 p95LatencyNs=583140 p99LatencyNs=921344 p999LatencyNs=2560325 maxLatencyNs=4205441 rssBytes=21184512 heapAllocBytes=6017800 heapSysBytes=14614528 heapObjects=81686 gcCount=13 gcPauseNs=2508668 goroutines=5 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.499419316s throughput=213415.95 ops/s
BenchmarkNetpollHTTP1-8 320000 4686 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 1.81739718s | 1 | 1 |  |
| 2 | 0 | 1.801193325s | 1 | 1 |  |
| 3 | 0 | 1.766149454s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netpoll | http1 |  | poller | - | 1024 | 64 | 5000 | 320000 | 0 | 204169.58 | 256833 | 1177092 | 24276992 | 13 | 1.567324571s |
| netpoll | http1 |  | poller | - | 1024 | 64 | 5000 | 320000 | 0 | 205662.31 | 259186 | 1038564 | 23552000 | 13 | 1.555948695s |
| netpoll | http1 |  | poller | - | 1024 | 64 | 5000 | 320000 | 0 | 213415.95 | 252640 | 921344 | 21184512 | 13 | 1.499419316s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNetpollHTTP1-8 | 320000 | 4898 | 0 | 0 |
| BenchmarkNetpollHTTP1-8 | 320000 | 4862 | 0 | 0 |
| BenchmarkNetpollHTTP1-8 | 320000 | 4686 | 0 | 0 |

### gnalloy epoll https1 tls11 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 45.128381907s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.1 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
```

Output:

```text
framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=662038 p95LatencyNs=1771667 p99LatencyNs=3350222 p999LatencyNs=5519387 maxLatencyNs=11061587 rssBytes=33824768 heapAllocBytes=10337736 heapSysBytes=23101440 heapObjects=87762 gcCount=79 gcPauseNs=16492009 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=4.046263488s throughput=79085.31 ops/s
BenchmarkGnalloyHTTPS1-8 320000 12645 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=660211 p95LatencyNs=1764203 p99LatencyNs=3269961 p999LatencyNs=4580806 maxLatencyNs=9267398 rssBytes=34107392 heapAllocBytes=8393920 heapSysBytes=22970368 heapObjects=48266 gcCount=78 gcPauseNs=19854888 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=4.003481011s throughput=79930.44 ops/s
BenchmarkGnalloyHTTPS1-8 320000 12511 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=638149 p95LatencyNs=1402810 p99LatencyNs=2384097 p999LatencyNs=4232801 maxLatencyNs=4864184 rssBytes=34328576 heapAllocBytes=7060856 heapSysBytes=23068672 heapObjects=21706 gcCount=79 gcPauseNs=21472665 goroutines=128 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=3.545805839s throughput=90247.47 ops/s
BenchmarkGnalloyHTTPS1-8 320000 11081 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 11.455955021s | 1 | 1 |  |
| 2 | 0 | 11.110980106s | 1 | 1 |  |
| 3 | 0 | 11.182465112s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 79085.31 | 662038 | 3350222 | 33824768 | 79 | 4.046263488s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 79930.44 | 660211 | 3269961 | 34107392 | 78 | 4.003481011s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 90247.47 | 638149 | 2384097 | 34328576 | 79 | 3.545805839s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 12645 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 12511 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 11081 | 0 | 0 |

### netty epoll https1 tls11 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 47.548118371s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https1 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn http/1.1 --tls-version 1.1 --cipher-suites TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
```

Output:

```text
framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=690883 p95LatencyNs=3348637 p99LatencyNs=5808241 p999LatencyNs=11222875 maxLatencyNs=15164945 rssBytes=398970880 heapAllocBytes=135362632 heapSysBytes=264241152 heapObjects=0 gcCount=13 gcPauseNs=33000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.588008089S throughput=57265.49 ops/s
BenchmarkNettyHTTPS1-8 320000 17463 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=613560 p95LatencyNs=3601636 p99LatencyNs=6444794 p999LatencyNs=11164067 maxLatencyNs=23942909 rssBytes=481230848 heapAllocBytes=215053880 heapSysBytes=318767104 heapObjects=0 gcCount=12 gcPauseNs=51000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.478327651S throughput=58411.99 ops/s
BenchmarkNettyHTTPS1-8 320000 17120 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=542254 p95LatencyNs=2668924 p99LatencyNs=4793318 p999LatencyNs=8709566 maxLatencyNs=12384420 rssBytes=408506368 heapAllocBytes=132539896 heapSysBytes=264241152 heapObjects=0 gcCount=13 gcPauseNs=21000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.568330367S throughput=70047.47 ops/s
BenchmarkNettyHTTPS1-8 320000 14276 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 11.858644685s | 1 | 1 |  |
| 2 | 0 | 12.080853028s | 1 | 1 |  |
| 3 | 0 | 11.165413332s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 57265.49 | 690883 | 5808241 | 398970880 | 13 | 5.588008089s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 58411.99 | 613560 | 6444794 | 481230848 | 12 | 5.478327651s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 70047.47 | 542254 | 4793318 | 408506368 | 13 | 4.568330367s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS1-8 | 320000 | 17463 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 17120 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 14276 | 0 | 0 |

### fasthttp https1 tls11 128B

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | https1 |
| backend | net |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 13.892251247s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol https1 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.1 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
```

Output:

```text
framework=fasthttp protocol=https1 backend=net tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=486025 p95LatencyNs=1215498 p99LatencyNs=2407510 p999LatencyNs=5138886 maxLatencyNs=9510153 rssBytes=22130688 heapAllocBytes=5557792 heapSysBytes=15007744 heapObjects=91467 gcCount=18 gcPauseNs=9127162 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.95587866s throughput=108258.84 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 9237 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=487409 p95LatencyNs=1203605 p99LatencyNs=2173329 p999LatencyNs=4327707 maxLatencyNs=16240291 rssBytes=23277568 heapAllocBytes=4734528 heapSysBytes=15073280 heapObjects=64747 gcCount=19 gcPauseNs=5334760 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.873795676s throughput=111350.99 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 8981 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=480240 p95LatencyNs=1083121 p99LatencyNs=1377070 p999LatencyNs=1807693 maxLatencyNs=2755570 rssBytes=22765568 heapAllocBytes=2923272 heapSysBytes=15237120 heapObjects=7361 gcCount=20 gcPauseNs=3468686 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.639480501s throughput=121235.98 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 8248 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 3.807480659s | 1 | 1 |  |
| 2 | 0 | 3.555394315s | 1 | 1 |  |
| 3 | 0 | 3.292380061s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 108258.84 | 486025 | 2407510 | 22130688 | 18 | 2.95587866s |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 111350.99 | 487409 | 2173329 | 23277568 | 19 | 2.873795676s |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 121235.98 | 480240 | 1377070 | 22765568 | 20 | 2.639480501s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 9237 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 8981 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 8248 | 0 | 0 |

### gnalloy epoll https1 tls11 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 49.152656914s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.1 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
```

Output:

```text
framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=716775 p95LatencyNs=1568291 p99LatencyNs=2772862 p999LatencyNs=5385032 maxLatencyNs=6707898 rssBytes=34344960 heapAllocBytes=8153752 heapSysBytes=23265280 heapObjects=39467 gcCount=77 gcPauseNs=11924335 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.962084652s throughput=80765.56 ops/s
BenchmarkGnalloyHTTPS1-8 320000 12382 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=745911 p95LatencyNs=2230996 p99LatencyNs=4384733 p999LatencyNs=6403626 maxLatencyNs=8125215 rssBytes=34045952 heapAllocBytes=10112472 heapSysBytes=23232512 heapObjects=79627 gcCount=75 gcPauseNs=18246979 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=4.715568535s throughput=67860.32 ops/s
BenchmarkGnalloyHTTPS1-8 320000 14736 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=715574 p95LatencyNs=1579086 p99LatencyNs=2483285 p999LatencyNs=4613909 maxLatencyNs=12114804 rssBytes=34226176 heapAllocBytes=10088120 heapSysBytes=23232512 heapObjects=79208 gcCount=76 gcPauseNs=12994108 goroutines=130 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.957339665s throughput=80862.40 ops/s
BenchmarkGnalloyHTTPS1-8 320000 12367 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 11.673556919s | 1 | 1 |  |
| 2 | 0 | 12.391914787s | 1 | 1 |  |
| 3 | 0 | 11.850374518s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 80765.56 | 716775 | 2772862 | 34344960 | 77 | 3.962084652s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 67860.32 | 745911 | 4384733 | 34045952 | 75 | 4.715568535s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 80862.4 | 715574 | 2483285 | 34226176 | 76 | 3.957339665s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 12382 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 14736 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 12367 | 0 | 0 |

### netty epoll https1 tls11 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 49.790517748s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https1 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn http/1.1 --tls-version 1.1 --cipher-suites TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
```

Output:

```text
framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=721893 p95LatencyNs=3558077 p99LatencyNs=6388434 p999LatencyNs=10445363 maxLatencyNs=12792727 rssBytes=424497152 heapAllocBytes=89546384 heapSysBytes=318767104 heapObjects=0 gcCount=18 gcPauseNs=40000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.774266216S throughput=55418.30 ops/s
BenchmarkNettyHTTPS1-8 320000 18045 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=694979 p95LatencyNs=3551809 p99LatencyNs=6314190 p999LatencyNs=12908420 maxLatencyNs=15918927 rssBytes=472879104 heapAllocBytes=181496496 heapSysBytes=318767104 heapObjects=0 gcCount=16 gcPauseNs=52000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.636557434S throughput=56772.24 ops/s
BenchmarkNettyHTTPS1-8 320000 17614 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=709149 p95LatencyNs=3016990 p99LatencyNs=5426107 p999LatencyNs=9407785 maxLatencyNs=11815380 rssBytes=407826432 heapAllocBytes=96545576 heapSysBytes=264241152 heapObjects=0 gcCount=18 gcPauseNs=37000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.388966547S throughput=59380.59 ops/s
BenchmarkNettyHTTPS1-8 320000 16841 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 12.448102788s | 1 | 1 |  |
| 2 | 0 | 12.904885351s | 1 | 1 |  |
| 3 | 0 | 12.187884003s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 55418.3 | 721893 | 6388434 | 424497152 | 18 | 5.774266216s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 56772.24 | 694979 | 6314190 | 472879104 | 16 | 5.636557434s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 59380.59 | 709149 | 5426107 | 407826432 | 18 | 5.388966547s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS1-8 | 320000 | 18045 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 17614 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 16841 | 0 | 0 |

### fasthttp https1 tls11 1KiB

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | https1 |
| backend | net |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 16.366035837s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol https1 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.1 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
```

Output:

```text
framework=fasthttp protocol=https1 backend=net tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=583914 p95LatencyNs=1599819 p99LatencyNs=3747313 p999LatencyNs=6432531 maxLatencyNs=10178777 rssBytes=23191552 heapAllocBytes=4987192 heapSysBytes=15073280 heapObjects=66579 gcCount=18 gcPauseNs=8672401 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.686516053s throughput=86802.82 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 11520 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=543298 p95LatencyNs=1254539 p99LatencyNs=1647843 p999LatencyNs=2291276 maxLatencyNs=4609454 rssBytes=20987904 heapAllocBytes=3143696 heapSysBytes=15269888 heapObjects=7996 gcCount=19 gcPauseNs=2940399 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.009131743s throughput=106342.97 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 9404 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.1 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=578274 p95LatencyNs=1537418 p99LatencyNs=3580365 p999LatencyNs=6223537 maxLatencyNs=12201025 rssBytes=24707072 heapAllocBytes=4862624 heapSysBytes=19300352 heapObjects=62950 gcCount=18 gcPauseNs=7985858 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.568093703s throughput=89683.74 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 11150 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 4.376958852s | 1 | 1 |  |
| 2 | 0 | 3.816500188s | 1 | 1 |  |
| 3 | 0 | 4.318090737s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 86802.82 | 583914 | 3747313 | 23191552 | 18 | 3.686516053s |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 106342.97 | 543298 | 1647843 | 20987904 | 19 | 3.009131743s |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 89683.74 | 578274 | 3580365 | 24707072 | 18 | 3.568093703s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 11520 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 9404 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 11150 | 0 | 0 |

### gnalloy epoll https1 tls12 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 38.347499506s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.2 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=479303 p95LatencyNs=1386869 p99LatencyNs=2676288 p999LatencyNs=4578237 maxLatencyNs=8350510 rssBytes=33828864 heapAllocBytes=9836344 heapSysBytes=23035904 heapObjects=72118 gcCount=73 gcPauseNs=14894737 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.967080049s throughput=107850.14 ops/s
BenchmarkGnalloyHTTPS1-8 320000 9272 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=469176 p95LatencyNs=860771 p99LatencyNs=1461935 p999LatencyNs=3471074 maxLatencyNs=5047073 rssBytes=31649792 heapAllocBytes=12896800 heapSysBytes=18907136 heapObjects=129806 gcCount=69 gcPauseNs=15238656 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.479131545s throughput=129077.46 ops/s
BenchmarkGnalloyHTTPS1-8 320000 7747 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=465889 p95LatencyNs=891157 p99LatencyNs=1351212 p999LatencyNs=3110582 maxLatencyNs=5063479 rssBytes=33808384 heapAllocBytes=10894208 heapSysBytes=23134208 heapObjects=91952 gcCount=67 gcPauseNs=9390693 goroutines=130 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.48529852s throughput=128757.17 ops/s
BenchmarkGnalloyHTTPS1-8 320000 7767 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 10.611904123s | 1 | 1 |  |
| 2 | 0 | 9.587956928s | 1 | 1 |  |
| 3 | 0 | 8.866533345s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 107850.14 | 479303 | 2676288 | 33828864 | 73 | 2.967080049s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 129077.46 | 469176 | 1461935 | 31649792 | 69 | 2.479131545s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 128757.17 | 465889 | 1351212 | 33808384 | 67 | 2.48529852s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 9272 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 7747 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 7767 | 0 | 0 |

### netty epoll https1 tls12 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 45.117450523s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https1 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn http/1.1 --tls-version 1.2 --cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=620835 p95LatencyNs=2782781 p99LatencyNs=5648456 p999LatencyNs=10740072 maxLatencyNs=16786734 rssBytes=475123712 heapAllocBytes=151752088 heapSysBytes=318767104 heapObjects=0 gcCount=17 gcPauseNs=46000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.980872844S throughput=64245.77 ops/s
BenchmarkNettyHTTPS1-8 320000 15565 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=466628 p95LatencyNs=3066467 p99LatencyNs=6198546 p999LatencyNs=11342986 maxLatencyNs=22833161 rssBytes=454008832 heapAllocBytes=197769736 heapSysBytes=318767104 heapObjects=0 gcCount=17 gcPauseNs=45000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.85721631S throughput=65881.36 ops/s
BenchmarkNettyHTTPS1-8 320000 15179 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=449524 p95LatencyNs=3067938 p99LatencyNs=6735373 p999LatencyNs=14705779 maxLatencyNs=26501753 rssBytes=605855744 heapAllocBytes=276251656 heapSysBytes=465567744 heapObjects=0 gcCount=14 gcPauseNs=97000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.753611882S throughput=67317.23 ops/s
BenchmarkNettyHTTPS1-8 320000 14855 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 11.466677075s | 1 | 1 |  |
| 2 | 0 | 11.059352351s | 1 | 1 |  |
| 3 | 0 | 11.068040724s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 64245.77 | 620835 | 5648456 | 475123712 | 17 | 4.980872844s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 65881.36 | 466628 | 6198546 | 454008832 | 17 | 4.85721631s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 67317.23 | 449524 | 6735373 | 605855744 | 14 | 4.753611882s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS1-8 | 320000 | 15565 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 15179 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 14855 | 0 | 0 |

### fasthttp https1 tls12 128B

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | https1 |
| backend | net |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 10.566773725s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol https1 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.2 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=fasthttp protocol=https1 backend=net tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=351086 p95LatencyNs=926365 p99LatencyNs=2716256 p999LatencyNs=6060532 maxLatencyNs=16305175 rssBytes=22441984 heapAllocBytes=5462296 heapSysBytes=10682368 heapObjects=158146 gcCount=3 gcPauseNs=658506 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.374899296s throughput=134742.56 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 7422 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=334937 p95LatencyNs=821582 p99LatencyNs=1622029 p999LatencyNs=4340097 maxLatencyNs=6337236 rssBytes=22396928 heapAllocBytes=3361416 heapSysBytes=14909440 heapObjects=28067 gcCount=4 gcPauseNs=601754 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.995535596s throughput=160357.95 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 6236 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=313726 p95LatencyNs=707331 p99LatencyNs=887788 p999LatencyNs=1177833 maxLatencyNs=14594279 rssBytes=22245376 heapAllocBytes=4480472 heapSysBytes=15106048 heapObjects=99528 gcCount=3 gcPauseNs=399260 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.730690289s throughput=184897.32 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 5408 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.995121321s | 1 | 1 |  |
| 2 | 0 | 2.603569747s | 1 | 1 |  |
| 3 | 0 | 2.377210499s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 134742.56 | 351086 | 2716256 | 22441984 | 3 | 2.374899296s |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 160357.95 | 334937 | 1622029 | 22396928 | 4 | 1.995535596s |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 184897.32 | 313726 | 887788 | 22245376 | 3 | 1.730690289s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 7422 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 6236 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 5408 | 0 | 0 |

### gnalloy epoll https1 tls12 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 38.111495524s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.2 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=493927 p95LatencyNs=1315765 p99LatencyNs=2971700 p999LatencyNs=5114367 maxLatencyNs=18951315 rssBytes=34320384 heapAllocBytes=8823512 heapSysBytes=23068672 heapObjects=49236 gcCount=68 gcPauseNs=19321536 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.023740431s throughput=105829.19 ops/s
BenchmarkGnalloyHTTPS1-8 320000 9449 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=494256 p95LatencyNs=1309644 p99LatencyNs=2613490 p999LatencyNs=5288928 maxLatencyNs=7031123 rssBytes=34476032 heapAllocBytes=9244824 heapSysBytes=23166976 heapObjects=57552 gcCount=67 gcPauseNs=18992839 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.073860152s throughput=104103.63 ops/s
BenchmarkGnalloyHTTPS1-8 320000 9606 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=489309 p95LatencyNs=1326446 p99LatencyNs=2869251 p999LatencyNs=4890745 maxLatencyNs=5972199 rssBytes=33828864 heapAllocBytes=7861808 heapSysBytes=23068672 heapObjects=31397 gcCount=68 gcPauseNs=25889842 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.051313597s throughput=104872.87 ops/s
BenchmarkGnalloyHTTPS1-8 320000 9535 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 10.02579253s | 1 | 1 |  |
| 2 | 0 | 9.062601715s | 1 | 1 |  |
| 3 | 0 | 9.321064579s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 105829.19 | 493927 | 2971700 | 34320384 | 68 | 3.023740431s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 104103.63 | 494256 | 2613490 | 34476032 | 67 | 3.073860152s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 104872.87 | 489309 | 2869251 | 33828864 | 68 | 3.051313597s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 9449 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 9606 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 9535 | 0 | 0 |

### netty epoll https1 tls12 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 46.55765398s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https1 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn http/1.1 --tls-version 1.2 --cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=482338 p95LatencyNs=3141491 p99LatencyNs=6048609 p999LatencyNs=12344123 maxLatencyNs=17473898 rssBytes=469520384 heapAllocBytes=216917888 heapSysBytes=318767104 heapObjects=0 gcCount=17 gcPauseNs=53000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.864282077S throughput=65785.66 ops/s
BenchmarkNettyHTTPS1-8 320000 15201 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=622821 p95LatencyNs=3079694 p99LatencyNs=6650647 p999LatencyNs=11992374 maxLatencyNs=24370202 rssBytes=557064192 heapAllocBytes=195901856 heapSysBytes=383778816 heapObjects=0 gcCount=16 gcPauseNs=67000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.205568357S throughput=61472.63 ops/s
BenchmarkNettyHTTPS1-8 320000 16267 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=583656 p95LatencyNs=3121778 p99LatencyNs=5931033 p999LatencyNs=12435492 maxLatencyNs=18859758 rssBytes=556355584 heapAllocBytes=154090128 heapSysBytes=383778816 heapObjects=0 gcCount=16 gcPauseNs=66000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.071903118S throughput=63092.69 ops/s
BenchmarkNettyHTTPS1-8 320000 15850 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 11.217252085s | 1 | 1 |  |
| 2 | 0 | 12.100093334s | 1 | 1 |  |
| 3 | 0 | 11.776223912s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 65785.66 | 482338 | 6048609 | 469520384 | 17 | 4.864282077s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 61472.63 | 622821 | 6650647 | 557064192 | 16 | 5.205568357s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 63092.69 | 583656 | 5931033 | 556355584 | 16 | 5.071903118s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS1-8 | 320000 | 15201 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 16267 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 15850 | 0 | 0 |

### fasthttp https1 tls12 1KiB

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | https1 |
| backend | net |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 10.21644126s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol https1 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.2 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=fasthttp protocol=https1 backend=net tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=343384 p95LatencyNs=819254 p99LatencyNs=1563392 p999LatencyNs=3909024 maxLatencyNs=6418329 rssBytes=20262912 heapAllocBytes=5397104 heapSysBytes=10682368 heapObjects=147242 gcCount=3 gcPauseNs=384514 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.047859036s throughput=156260.76 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 6400 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=350663 p95LatencyNs=872972 p99LatencyNs=2115709 p999LatencyNs=4419992 maxLatencyNs=6151505 rssBytes=22208512 heapAllocBytes=5415064 heapSysBytes=15007744 heapObjects=149169 gcCount=3 gcPauseNs=666821 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.071544905s throughput=154474.08 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 6474 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=325257 p95LatencyNs=729007 p99LatencyNs=919472 p999LatencyNs=1246958 maxLatencyNs=2761144 rssBytes=21839872 heapAllocBytes=4208104 heapSysBytes=15073280 heapObjects=75169 gcCount=4 gcPauseNs=676505 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.789032322s throughput=178867.65 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 5591 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.670076246s | 1 | 1 |  |
| 2 | 0 | 2.611361897s | 1 | 1 |  |
| 3 | 0 | 2.415332401s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 156260.76 | 343384 | 1563392 | 20262912 | 3 | 2.047859036s |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 154474.08 | 350663 | 2115709 | 22208512 | 3 | 2.071544905s |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 178867.65 | 325257 | 919472 | 21839872 | 4 | 1.789032322s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 6400 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 6474 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 5591 | 0 | 0 |

### gnalloy epoll https1 tls13 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 1m31.89304899s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=484186 p95LatencyNs=1295680 p99LatencyNs=3067208 p999LatencyNs=5384925 maxLatencyNs=5989485 rssBytes=33853440 heapAllocBytes=9891744 heapSysBytes=23003136 heapObjects=76054 gcCount=61 gcPauseNs=15693790 goroutines=128 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.894483426s throughput=110555.13 ops/s
BenchmarkGnalloyHTTPS1-8 320000 9045 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=479892 p95LatencyNs=1253882 p99LatencyNs=2762250 p999LatencyNs=5390842 maxLatencyNs=5774146 rssBytes=33980416 heapAllocBytes=9326808 heapSysBytes=23003136 heapObjects=65251 gcCount=62 gcPauseNs=18190684 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.955328841s throughput=108278.98 ops/s
BenchmarkGnalloyHTTPS1-8 320000 9235 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=495476 p95LatencyNs=1665174 p99LatencyNs=3464141 p999LatencyNs=6787371 maxLatencyNs=8357052 rssBytes=33755136 heapAllocBytes=6668520 heapSysBytes=23035904 heapObjects=15429 gcCount=62 gcPauseNs=20326202 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=3.363933002s throughput=95126.75 ops/s
BenchmarkGnalloyHTTPS1-8 320000 10512 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 22.938635379s | 1 | 1 |  |
| 2 | 0 | 22.96778413s | 1 | 1 |  |
| 3 | 0 | 23.371243163s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 110555.13 | 484186 | 3067208 | 33853440 | 61 | 2.894483426s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 108278.98 | 479892 | 2762250 | 33980416 | 62 | 2.955328841s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 95126.75 | 495476 | 3464141 | 33755136 | 62 | 3.363933002s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 9045 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 9235 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 10512 | 0 | 0 |

### netty epoll https1 tls13 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https1 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 49.125129318s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https1 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn http/1.1 --tls-version 1.3
```

Output:

```text
framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=624711 p95LatencyNs=3317066 p99LatencyNs=7709098 p999LatencyNs=15860370 maxLatencyNs=23742458 rssBytes=638709760 heapAllocBytes=183597072 heapSysBytes=461373440 heapObjects=0 gcCount=14 gcPauseNs=95000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.902199742S throughput=54217.07 ops/s
BenchmarkNettyHTTPS1-8 320000 18444 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=551238 p95LatencyNs=2626127 p99LatencyNs=5082296 p999LatencyNs=11066496 maxLatencyNs=15980585 rssBytes=638586880 heapAllocBytes=74619656 heapSysBytes=461373440 heapObjects=0 gcCount=15 gcPauseNs=73000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.724081562S throughput=67738.03 ops/s
BenchmarkNettyHTTPS1-8 320000 14763 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=596093 p95LatencyNs=3067250 p99LatencyNs=6375818 p999LatencyNs=14184606 maxLatencyNs=17374272 rssBytes=495542272 heapAllocBytes=82336040 heapSysBytes=318767104 heapObjects=0 gcCount=17 gcPauseNs=64000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.542090962S throughput=57739.94 ops/s
BenchmarkNettyHTTPS1-8 320000 17319 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 12.592302624s | 1 | 1 |  |
| 2 | 0 | 12.008465526s | 1 | 1 |  |
| 3 | 0 | 12.286380992s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 54217.07 | 624711 | 7709098 | 638709760 | 14 | 5.902199742s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 67738.03 | 551238 | 5082296 | 638586880 | 15 | 4.724081562s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 57739.94 | 596093 | 6375818 | 495542272 | 17 | 5.542090962s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS1-8 | 320000 | 18444 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 14763 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 17319 | 0 | 0 |

### fasthttp https1 tls13 128B

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | https1 |
| backend | net |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 9.567843496s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol https1 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.3
```

Output:

```text
framework=fasthttp protocol=https1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=337189 p95LatencyNs=875742 p99LatencyNs=1914047 p999LatencyNs=4364391 maxLatencyNs=13492840 rssBytes=24629248 heapAllocBytes=3508432 heapSysBytes=14876672 heapObjects=23630 gcCount=9 gcPauseNs=897612 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.077427067s throughput=154036.70 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 6492 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=331402 p95LatencyNs=775207 p99LatencyNs=1365796 p999LatencyNs=3840707 maxLatencyNs=6173233 rssBytes=22573056 heapAllocBytes=5008184 heapSysBytes=14876672 heapObjects=117761 gcCount=7 gcPauseNs=558769 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.959175808s throughput=163333.99 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 6122 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=312970 p95LatencyNs=705784 p99LatencyNs=942106 p999LatencyNs=1367142 maxLatencyNs=4751551 rssBytes=24563712 heapAllocBytes=3793664 heapSysBytes=15073280 heapObjects=43938 gcCount=8 gcPauseNs=1343679 goroutines=67 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=1.75241862s throughput=182604.77 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 5476 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.600193673s | 1 | 1 |  |
| 2 | 0 | 2.332405874s | 1 | 1 |  |
| 3 | 0 | 2.329097964s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 154036.7 | 337189 | 1914047 | 24629248 | 9 | 2.077427067s |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 163333.99 | 331402 | 1365796 | 22573056 | 7 | 1.959175808s |
| fasthttp | https1 | http/1.1 | net | - | 128 | 64 | 5000 | 320000 | 0 | 182604.77 | 312970 | 942106 | 24563712 | 8 | 1.75241862s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 6492 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 6122 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 5476 | 0 | 0 |

### gnalloy epoll https1 tls13 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 1m32.435656998s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https1 -backend epoll -http1-mode raw -read-buffer-size 384 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=477991 p95LatencyNs=1130500 p99LatencyNs=2579287 p999LatencyNs=4823060 maxLatencyNs=9732924 rssBytes=34033664 heapAllocBytes=7635312 heapSysBytes=23068672 heapObjects=30070 gcCount=60 gcPauseNs=19209063 goroutines=126 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.770554536s throughput=115500.34 ops/s
BenchmarkGnalloyHTTPS1-8 320000 8658 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=516744 p95LatencyNs=1940231 p99LatencyNs=4412951 p999LatencyNs=6588304 maxLatencyNs=8964895 rssBytes=32002048 heapAllocBytes=11308888 heapSysBytes=18776064 heapObjects=98522 gcCount=59 gcPauseNs=21525403 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.673374399s throughput=87113.36 ops/s
BenchmarkGnalloyHTTPS1-8 320000 11479 ns/op

framework=gnalloy protocol=https1 backend=epoll http1Mode=raw boss=1 workers=4 readBufferSize=384 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=492612 p95LatencyNs=1345017 p99LatencyNs=2652940 p999LatencyNs=5471721 maxLatencyNs=6400653 rssBytes=34029568 heapAllocBytes=8395336 heapSysBytes=23035904 heapObjects=43965 gcCount=60 gcPauseNs=18989887 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.057315433s throughput=104666.99 ops/s
BenchmarkGnalloyHTTPS1-8 320000 9554 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 22.831264193s | 1 | 1 |  |
| 2 | 0 | 23.788283266s | 1 | 1 |  |
| 3 | 0 | 23.10154336s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 115500.34 | 477991 | 2579287 | 34033664 | 60 | 2.770554536s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 87113.36 | 516744 | 4412951 | 32002048 | 59 | 3.673374399s |
| gnalloy | https1 | http/1.1 | epoll | boss=1 workers=4 readBuffer=384 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 104666.99 | 492612 | 2652940 | 34029568 | 60 | 3.057315433s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 8658 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 11479 | 0 | 0 |
| BenchmarkGnalloyHTTPS1-8 | 320000 | 9554 | 0 | 0 |

### netty epoll https1 tls13 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https1 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 52.337908368s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https1 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn http/1.1 --tls-version 1.3
```

Output:

```text
framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=622365 p95LatencyNs=4047604 p99LatencyNs=8248798 p999LatencyNs=12747886 maxLatencyNs=16189512 rssBytes=489738240 heapAllocBytes=179860032 heapSysBytes=318767104 heapObjects=0 gcCount=19 gcPauseNs=56000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.995343156S throughput=53374.76 ops/s
BenchmarkNettyHTTPS1-8 320000 18735 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=954166 p95LatencyNs=3447233 p99LatencyNs=6595478 p999LatencyNs=10928943 maxLatencyNs=16887993 rssBytes=572129280 heapAllocBytes=93532496 heapSysBytes=383778816 heapObjects=0 gcCount=17 gcPauseNs=81000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT6.677515907S throughput=47922.01 ops/s
BenchmarkNettyHTTPS1-8 320000 20867 ns/op

framework=netty protocol=https1 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=591124 p95LatencyNs=3717257 p99LatencyNs=6915962 p999LatencyNs=14736634 maxLatencyNs=22838169 rssBytes=642633728 heapAllocBytes=302580664 heapSysBytes=461373440 heapObjects=0 gcCount=15 gcPauseNs=94000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.785935793S throughput=55306.52 ops/s
BenchmarkNettyHTTPS1-8 320000 18081 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 12.648729799s | 1 | 1 |  |
| 2 | 0 | 13.931830714s | 1 | 1 |  |
| 3 | 0 | 13.012902476s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 53374.76 | 622365 | 8248798 | 489738240 | 19 | 5.995343156s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 47922.01 | 954166 | 6595478 | 572129280 | 17 | 6.677515907s |
| netty | https1 | http/1.1 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 55306.52 | 591124 | 6915962 | 642633728 | 15 | 5.785935793s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS1-8 | 320000 | 18735 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 20867 | 0 | 0 |
| BenchmarkNettyHTTPS1-8 | 320000 | 18081 | 0 | 0 |

### fasthttp https1 tls13 1KiB

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | https1 |
| backend | net |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 9.661590915s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/fasthttp-bench -protocol https1 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn http/1.1 -tls-version 1.3
```

Output:

```text
framework=fasthttp protocol=https1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=320449 p95LatencyNs=739983 p99LatencyNs=989190 p999LatencyNs=1236408 maxLatencyNs=2406334 rssBytes=24592384 heapAllocBytes=4273752 heapSysBytes=15073280 heapObjects=66378 gcCount=7 gcPauseNs=629820 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.830367854s throughput=174828.25 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 5720 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=321359 p95LatencyNs=742429 p99LatencyNs=949962 p999LatencyNs=1253169 maxLatencyNs=4045509 rssBytes=22523904 heapAllocBytes=5317040 heapSysBytes=10911744 heapObjects=131262 gcCount=7 gcPauseNs=687464 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.811576454s throughput=176641.73 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 5661 ns/op

framework=fasthttp protocol=https1 backend=net tlsVersion=1.3 cipherSuites= negotiatedProtocol=http/1.1 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=343144 p95LatencyNs=813466 p99LatencyNs=1235495 p999LatencyNs=3126126 maxLatencyNs=5597617 rssBytes=22499328 heapAllocBytes=6036256 heapSysBytes=15007744 heapObjects=175411 gcCount=7 gcPauseNs=838763 goroutines=67 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=1.981127928s throughput=161524.15 ops/s
BenchmarkFastHTTPHTTPS1-8 320000 6191 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 2.230096934s | 1 | 1 |  |
| 2 | 0 | 2.379297368s | 1 | 1 |  |
| 3 | 0 | 2.498676562s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 174828.25 | 320449 | 989190 | 24592384 | 7 | 1.830367854s |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 176641.73 | 321359 | 949962 | 22523904 | 7 | 1.811576454s |
| fasthttp | https1 | http/1.1 | net | - | 1024 | 64 | 5000 | 320000 | 0 | 161524.15 | 343144 | 1235495 | 22499328 | 7 | 1.981127928s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 5720 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 5661 | 0 | 0 |
| BenchmarkFastHTTPHTTPS1-8 | 320000 | 6191 | 0 | 0 |

### gnet https1 unsupported

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | https1 |
| backend | unsupported |
| payload | - |
| duration | 118ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
gnet benchmark harness intentionally measures the native poller and does not wrap TLS
```

### netpoll https1 unsupported

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | https1 |
| backend | unsupported |
| payload | - |
| duration | 94ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
netpoll benchmark harness intentionally measures the native poller and does not wrap TLS
```

### gnalloy epoll http2 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | http2 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 14.182011317s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol http2 -backend epoll -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=http2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=514297 p95LatencyNs=1098814 p99LatencyNs=1710540 p999LatencyNs=3971895 maxLatencyNs=7311551 rssBytes=22220800 heapAllocBytes=7206384 heapSysBytes=15728640 heapObjects=87120 gcCount=118 gcPauseNs=16606489 goroutines=9 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=2.877561691s throughput=111205.26 ops/s
BenchmarkGnalloyHTTP2-8 320000 8992 ns/op

framework=gnalloy protocol=http2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=522503 p95LatencyNs=1460957 p99LatencyNs=2717045 p999LatencyNs=5087456 maxLatencyNs=9704671 rssBytes=23523328 heapAllocBytes=4807888 heapSysBytes=15630336 heapObjects=29120 gcCount=116 gcPauseNs=31212034 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=3.29550295s throughput=97102.02 ops/s
BenchmarkGnalloyHTTP2-8 320000 10298 ns/op

framework=gnalloy protocol=http2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=514105 p95LatencyNs=1296027 p99LatencyNs=2248865 p999LatencyNs=4098438 maxLatencyNs=12267880 rssBytes=24104960 heapAllocBytes=7314464 heapSysBytes=15695872 heapObjects=89387 gcCount=115 gcPauseNs=23404855 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=3.075362497s throughput=104052.77 ops/s
BenchmarkGnalloyHTTP2-8 320000 9611 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 3.264918376s | 1 | 1 |  |
| 2 | 0 | 3.762794065s | 1 | 1 |  |
| 3 | 0 | 3.462372422s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 111205.26 | 514297 | 1710540 | 22220800 | 118 | 2.877561691s |
| gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 97102.02 | 522503 | 2717045 | 23523328 | 116 | 3.29550295s |
| gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 104052.77 | 514105 | 2248865 | 24104960 | 115 | 3.075362497s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTP2-8 | 320000 | 8992 | 0 | 0 |
| BenchmarkGnalloyHTTP2-8 | 320000 | 10298 | 0 | 0 |
| BenchmarkGnalloyHTTP2-8 | 320000 | 9611 | 0 | 0 |

### netty epoll http2 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | http2 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 39.552832531s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol http2 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=http2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=572205 p95LatencyNs=2941290 p99LatencyNs=5944494 p999LatencyNs=9633166 maxLatencyNs=16627873 rssBytes=318230528 heapAllocBytes=148644152 heapSysBytes=264241152 heapObjects=0 gcCount=5 gcPauseNs=28000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.946943855S throughput=64686.40 ops/s
BenchmarkNettyHTTP2-8 320000 15459 ns/op

framework=netty protocol=http2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=560288 p95LatencyNs=3524804 p99LatencyNs=7161116 p999LatencyNs=10856356 maxLatencyNs=13813003 rssBytes=322060288 heapAllocBytes=131618456 heapSysBytes=264241152 heapObjects=0 gcCount=5 gcPauseNs=19000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT5.572604866S throughput=57423.77 ops/s
BenchmarkNettyHTTP2-8 320000 17414 ns/op

framework=netty protocol=http2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=467540 p95LatencyNs=2535196 p99LatencyNs=5161723 p999LatencyNs=9123202 maxLatencyNs=13222239 rssBytes=314843136 heapAllocBytes=132258248 heapSysBytes=264241152 heapObjects=0 gcCount=5 gcPauseNs=14000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.425825611S throughput=72302.89 ops/s
BenchmarkNettyHTTP2-8 320000 13831 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 9.681688343s | 1 | 1 |  |
| 2 | 0 | 10.33015541s | 1 | 1 |  |
| 3 | 0 | 9.227028669s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | http2 |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 64686.4 | 572205 | 5944494 | 318230528 | 5 | 4.946943855s |
| netty | http2 |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 57423.77 | 560288 | 7161116 | 322060288 | 5 | 5.572604866s |
| netty | http2 |  | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 72302.89 | 467540 | 5161723 | 314843136 | 5 | 4.425825611s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTP2-8 | 320000 | 15459 | 0 | 0 |
| BenchmarkNettyHTTP2-8 | 320000 | 17414 | 0 | 0 |
| BenchmarkNettyHTTP2-8 | 320000 | 13831 | 0 | 0 |

### gnalloy epoll http2 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | http2 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 14.357215472s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol http2 -backend epoll -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m
```

Output:

```text
framework=gnalloy protocol=http2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=529573 p95LatencyNs=1518973 p99LatencyNs=2883792 p999LatencyNs=5458472 maxLatencyNs=8848723 rssBytes=23633920 heapAllocBytes=5371040 heapSysBytes=15630336 heapObjects=38838 gcCount=109 gcPauseNs=24461662 goroutines=9 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.424251081s throughput=93451.09 ops/s
BenchmarkGnalloyHTTP2-8 320000 10701 ns/op

framework=gnalloy protocol=http2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=517165 p95LatencyNs=1183473 p99LatencyNs=1723684 p999LatencyNs=3267686 maxLatencyNs=6259321 rssBytes=23724032 heapAllocBytes=5772936 heapSysBytes=15794176 heapObjects=48819 gcCount=112 gcPauseNs=18805502 goroutines=9 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=2.900759653s throughput=110315.93 ops/s
BenchmarkGnalloyHTTP2-8 320000 9065 ns/op

framework=gnalloy protocol=http2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=531091 p95LatencyNs=1473617 p99LatencyNs=2746669 p999LatencyNs=5544263 maxLatencyNs=7520613 rssBytes=23928832 heapAllocBytes=4338000 heapSysBytes=15630336 heapObjects=13944 gcCount=110 gcPauseNs=31880589 goroutines=8 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=3.307536665s throughput=96748.74 ops/s
BenchmarkGnalloyHTTP2-8 320000 10336 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 3.814015981s | 1 | 1 |  |
| 2 | 0 | 3.360266805s | 1 | 1 |  |
| 3 | 0 | 3.728057214s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 93451.09 | 529573 | 2883792 | 23633920 | 109 | 3.424251081s |
| gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 110315.93 | 517165 | 1723684 | 23724032 | 112 | 2.900759653s |
| gnalloy | http2 |  | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 96748.74 | 531091 | 2746669 | 23928832 | 110 | 3.307536665s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTP2-8 | 320000 | 10701 | 0 | 0 |
| BenchmarkGnalloyHTTP2-8 | 320000 | 9065 | 0 | 0 |
| BenchmarkGnalloyHTTP2-8 | 320000 | 10336 | 0 | 0 |

### netty epoll http2 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | http2 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 37.654233774s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol http2 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m
```

Output:

```text
framework=netty protocol=http2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=617661 p95LatencyNs=2812814 p99LatencyNs=5452168 p999LatencyNs=11375889 maxLatencyNs=22875258 rssBytes=331104256 heapAllocBytes=136875496 heapSysBytes=264241152 heapObjects=0 gcCount=5 gcPauseNs=18000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.883571776S throughput=65525.81 ops/s
BenchmarkNettyHTTP2-8 320000 15261 ns/op

framework=netty protocol=http2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=560318 p95LatencyNs=2703696 p99LatencyNs=5242324 p999LatencyNs=10340276 maxLatencyNs=14141910 rssBytes=339124224 heapAllocBytes=117803472 heapSysBytes=264241152 heapObjects=0 gcCount=5 gcPauseNs=18000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.833335701S throughput=66206.86 ops/s
BenchmarkNettyHTTP2-8 320000 15104 ns/op

framework=netty protocol=http2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol= latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=534770 p95LatencyNs=2451734 p99LatencyNs=4590663 p999LatencyNs=9516926 maxLatencyNs=13443430 rssBytes=307015680 heapAllocBytes=115682824 heapSysBytes=264241152 heapObjects=0 gcCount=5 gcPauseNs=13000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT4.448934785S throughput=71927.33 ops/s
BenchmarkNettyHTTP2-8 320000 13903 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 9.750928658s | 1 | 1 |  |
| 2 | 0 | 9.660467109s | 1 | 1 |  |
| 3 | 0 | 9.236297218s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | http2 |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 65525.81 | 617661 | 5452168 | 331104256 | 5 | 4.883571776s |
| netty | http2 |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 66206.86 | 560318 | 5242324 | 339124224 | 5 | 4.833335701s |
| netty | http2 |  | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 71927.33 | 534770 | 4590663 | 307015680 | 5 | 4.448934785s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTP2-8 | 320000 | 15261 | 0 | 0 |
| BenchmarkNettyHTTP2-8 | 320000 | 15104 | 0 | 0 |
| BenchmarkNettyHTTP2-8 | 320000 | 13903 | 0 | 0 |

### gnet http2 unsupported

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | http2 |
| backend | unsupported |
| payload | - |
| duration | 154ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
gnet does not provide an HTTP/2 codec in this benchmark harness
```

### fasthttp http2 unsupported

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | http2 |
| backend | unsupported |
| payload | - |
| duration | 86ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
fasthttp does not provide an HTTP/2 server implementation
```

### netpoll http2 unsupported

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | http2 |
| backend | unsupported |
| payload | - |
| duration | 90ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
netpoll benchmark harness does not include an HTTP/2 codec
```

### gnalloy epoll https2 tls11 unsupported

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https2 |
| backend | unsupported |
| payload | - |
| duration | 82ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
HTTP/2 over TLS requires TLS 1.2 or newer
```

### netty epoll https2 tls11 unsupported

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https2 |
| backend | unsupported |
| payload | - |
| duration | 85ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
HTTP/2 over TLS requires TLS 1.2 or newer
```

### gnalloy epoll https2 tls12 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https2 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 48.989869482s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https2 -backend epoll -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn h2 -tls-version 1.2 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=904544 p95LatencyNs=2067297 p99LatencyNs=3722814 p999LatencyNs=7123604 maxLatencyNs=11951881 rssBytes=34066432 heapAllocBytes=11579440 heapSysBytes=23101440 heapObjects=106317 gcCount=158 gcPauseNs=40821547 goroutines=126 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=5.267071791s throughput=60754.82 ops/s
BenchmarkGnalloyHTTPS2-8 320000 16460 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=891580 p95LatencyNs=1749322 p99LatencyNs=2289368 p999LatencyNs=4907506 maxLatencyNs=8012873 rssBytes=34111488 heapAllocBytes=13040104 heapSysBytes=23134208 heapObjects=136911 gcCount=160 gcPauseNs=22556095 goroutines=128 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=4.804931579s throughput=66598.24 ops/s
BenchmarkGnalloyHTTPS2-8 320000 15015 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=913054 p95LatencyNs=2221082 p99LatencyNs=3751254 p999LatencyNs=5838891 maxLatencyNs=6899376 rssBytes=33533952 heapAllocBytes=12144160 heapSysBytes=23068672 heapObjects=118741 gcCount=162 gcPauseNs=40415108 goroutines=133 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=5.339317711s throughput=59932.75 ops/s
BenchmarkGnalloyHTTPS2-8 320000 16685 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 12.201006822s | 1 | 1 |  |
| 2 | 0 | 12.022138744s | 1 | 1 |  |
| 3 | 0 | 12.922747104s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 60754.82 | 904544 | 3722814 | 34066432 | 158 | 5.267071791s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 66598.24 | 891580 | 2289368 | 34111488 | 160 | 4.804931579s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 59932.75 | 913054 | 3751254 | 33533952 | 162 | 5.339317711s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 16460 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 15015 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 16685 | 0 | 0 |

### netty epoll https2 tls12 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https2 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 57.26431588s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https2 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn h2 --tls-version 1.2 --cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=939795 p95LatencyNs=3561144 p99LatencyNs=6998574 p999LatencyNs=12068162 maxLatencyNs=21302099 rssBytes=481566720 heapAllocBytes=98401928 heapSysBytes=318767104 heapObjects=0 gcCount=22 gcPauseNs=58000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT6.687973352S throughput=47847.08 ops/s
BenchmarkNettyHTTPS2-8 320000 20900 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=913984 p95LatencyNs=4206092 p99LatencyNs=7597433 p999LatencyNs=11736299 maxLatencyNs=20474936 rssBytes=515072000 heapAllocBytes=130647184 heapSysBytes=318767104 heapObjects=0 gcCount=22 gcPauseNs=80000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT7.362504917S throughput=43463.47 ops/s
BenchmarkNettyHTTPS2-8 320000 23008 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1007541 p95LatencyNs=3974116 p99LatencyNs=6770355 p999LatencyNs=13851434 maxLatencyNs=28420713 rssBytes=620199936 heapAllocBytes=289755632 heapSysBytes=461373440 heapObjects=0 gcCount=17 gcPauseNs=104000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT7.598915011S throughput=42111.28 ops/s
BenchmarkNettyHTTPS2-8 320000 23747 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 13.955790874s | 1 | 1 |  |
| 2 | 0 | 14.384529382s | 1 | 1 |  |
| 3 | 0 | 14.86846229s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https2 | h2 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 47847.08 | 939795 | 6998574 | 481566720 | 22 | 6.687973352s |
| netty | https2 | h2 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 43463.47 | 913984 | 7597433 | 515072000 | 22 | 7.362504917s |
| netty | https2 | h2 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 42111.28 | 1007541 | 6770355 | 620199936 | 17 | 7.598915011s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS2-8 | 320000 | 20900 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 23008 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 23747 | 0 | 0 |

### gnalloy epoll https2 tls12 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https2 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 49.749221117s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https2 -backend epoll -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn h2 -tls-version 1.2 -cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=919995 p95LatencyNs=2353745 p99LatencyNs=4183646 p999LatencyNs=5946803 maxLatencyNs=6687492 rssBytes=34041856 heapAllocBytes=13830920 heapSysBytes=23035904 heapObjects=149123 gcCount=153 gcPauseNs=49251952 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=5.569805781s throughput=57452.63 ops/s
BenchmarkGnalloyHTTPS2-8 320000 17406 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=904240 p95LatencyNs=1754255 p99LatencyNs=2354118 p999LatencyNs=3961579 maxLatencyNs=4936362 rssBytes=33648640 heapAllocBytes=12860024 heapSysBytes=23068672 heapObjects=128025 gcCount=154 gcPauseNs=23903780 goroutines=128 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=4.859735623s throughput=65847.20 ops/s
BenchmarkGnalloyHTTPS2-8 320000 15187 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=922009 p95LatencyNs=2247585 p99LatencyNs=3762517 p999LatencyNs=6011498 maxLatencyNs=13101689 rssBytes=35962880 heapAllocBytes=11902016 heapSysBytes=23003136 heapObjects=107478 gcCount=155 gcPauseNs=52364212 goroutines=124 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=5.450552487s throughput=58709.64 ops/s
BenchmarkGnalloyHTTPS2-8 320000 17033 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 12.244186441s | 1 | 1 |  |
| 2 | 0 | 11.852220174s | 1 | 1 |  |
| 3 | 0 | 13.679031721s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 57452.63 | 919995 | 4183646 | 34041856 | 153 | 5.569805781s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 65847.2 | 904240 | 2354118 | 33648640 | 154 | 4.859735623s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 58709.64 | 922009 | 3762517 | 35962880 | 155 | 5.450552487s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 17406 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 15187 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 17033 | 0 | 0 |

### netty epoll https2 tls12 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https2 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 58.7732995s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https2 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn h2 --tls-version 1.2 --cipher-suites TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

Output:

```text
framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=857802 p95LatencyNs=3575897 p99LatencyNs=6364597 p999LatencyNs=12256564 maxLatencyNs=17303208 rssBytes=542945280 heapAllocBytes=206850232 heapSysBytes=383778816 heapObjects=0 gcCount=20 gcPauseNs=60000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT6.956806426S throughput=45998.12 ops/s
BenchmarkNettyHTTPS2-8 320000 21740 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=889439 p95LatencyNs=4175543 p99LatencyNs=8228167 p999LatencyNs=13304440 maxLatencyNs=25834600 rssBytes=419569664 heapAllocBytes=79613184 heapSysBytes=264241152 heapObjects=0 gcCount=25 gcPauseNs=62000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT7.325657767S throughput=43682.08 ops/s
BenchmarkNettyHTTPS2-8 320000 22893 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.2 cipherSuites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1002561 p95LatencyNs=4349300 p99LatencyNs=7733796 p999LatencyNs=12634489 maxLatencyNs=22735925 rssBytes=477458432 heapAllocBytes=206022672 heapSysBytes=318767104 heapObjects=0 gcCount=23 gcPauseNs=57000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT7.401366541S throughput=43235.26 ops/s
BenchmarkNettyHTTPS2-8 320000 23129 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 14.073204356s | 1 | 1 |  |
| 2 | 0 | 15.110246397s | 1 | 1 |  |
| 3 | 0 | 14.215022432s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https2 | h2 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 45998.12 | 857802 | 6364597 | 542945280 | 20 | 6.956806426s |
| netty | https2 | h2 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 43682.08 | 889439 | 8228167 | 419569664 | 25 | 7.325657767s |
| netty | https2 | h2 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 43235.26 | 1002561 | 7733796 | 477458432 | 23 | 7.401366541s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS2-8 | 320000 | 21740 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 22893 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 23129 | 0 | 0 |

### gnalloy epoll https2 tls13 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https2 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 24.385343687s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https2 -backend epoll -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn h2 -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=904533 p95LatencyNs=2131881 p99LatencyNs=3574613 p999LatencyNs=5837982 maxLatencyNs=9242466 rssBytes=34082816 heapAllocBytes=10265088 heapSysBytes=23035904 heapObjects=81993 gcCount=156 gcPauseNs=36771776 goroutines=128 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=5.272091772s throughput=60696.97 ops/s
BenchmarkGnalloyHTTPS2-8 320000 16475 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=890421 p95LatencyNs=1836211 p99LatencyNs=2826053 p999LatencyNs=5497679 maxLatencyNs=7182829 rssBytes=33566720 heapAllocBytes=9199144 heapSysBytes=23035904 heapObjects=60327 gcCount=156 gcPauseNs=33368516 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=4.887820103s throughput=65468.86 ops/s
BenchmarkGnalloyHTTPS2-8 320000 15274 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=905395 p95LatencyNs=2234598 p99LatencyNs=3819385 p999LatencyNs=5863365 maxLatencyNs=7414014 rssBytes=34009088 heapAllocBytes=10437200 heapSysBytes=22970368 heapObjects=85397 gcCount=155 gcPauseNs=49785318 goroutines=132 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=5.401311761s throughput=59244.87 ops/s
BenchmarkGnalloyHTTPS2-8 320000 16879 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 6.200900891s | 1 | 1 |  |
| 2 | 0 | 5.688232684s | 1 | 1 |  |
| 3 | 0 | 6.435858264s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 60696.97 | 904533 | 3574613 | 34082816 | 156 | 5.272091772s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 65468.86 | 890421 | 2826053 | 33566720 | 156 | 4.887820103s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 59244.87 | 905395 | 3819385 | 34009088 | 155 | 5.401311761s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 16475 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 15274 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 16879 | 0 | 0 |

### netty epoll https2 tls13 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https2 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 1m2.9600318s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https2 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn h2 --tls-version 1.3
```

Output:

```text
framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1240749 p95LatencyNs=4422166 p99LatencyNs=7778366 p999LatencyNs=12930616 maxLatencyNs=19827933 rssBytes=519393280 heapAllocBytes=166082632 heapSysBytes=318767104 heapObjects=0 gcCount=21 gcPauseNs=75000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT8.277371382S throughput=38659.62 ops/s
BenchmarkNettyHTTPS2-8 320000 25867 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1101032 p95LatencyNs=4204953 p99LatencyNs=8194155 p999LatencyNs=14317150 maxLatencyNs=24726827 rssBytes=664047616 heapAllocBytes=68944112 heapSysBytes=461373440 heapObjects=0 gcCount=18 gcPauseNs=103000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT7.860627915S throughput=40709.22 ops/s
BenchmarkNettyHTTPS2-8 320000 24564 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1110041 p95LatencyNs=4804009 p99LatencyNs=9108579 p999LatencyNs=16622497 maxLatencyNs=23313756 rssBytes=519905280 heapAllocBytes=68620056 heapSysBytes=318767104 heapObjects=0 gcCount=23 gcPauseNs=68000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT8.220864798S throughput=38925.35 ops/s
BenchmarkNettyHTTPS2-8 320000 25690 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 15.945876801s | 1 | 1 |  |
| 2 | 0 | 15.849344277s | 1 | 1 |  |
| 3 | 0 | 15.760603866s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https2 | h2 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 38659.62 | 1240749 | 7778366 | 519393280 | 21 | 8.277371382s |
| netty | https2 | h2 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 40709.22 | 1101032 | 8194155 | 664047616 | 18 | 7.860627915s |
| netty | https2 | h2 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 38925.35 | 1110041 | 9108579 | 519905280 | 23 | 8.220864797s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS2-8 | 320000 | 25867 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 24564 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 25690 | 0 | 0 |

### gnalloy epoll https2 tls13 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | https2 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 24.212259324s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol https2 -backend epoll -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn h2 -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=915754 p95LatencyNs=2156505 p99LatencyNs=3665159 p999LatencyNs=6316826 maxLatencyNs=9608169 rssBytes=33906688 heapAllocBytes=13434784 heapSysBytes=22970368 heapObjects=141763 gcCount=148 gcPauseNs=46116685 goroutines=132 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=5.297330579s throughput=60407.78 ops/s
BenchmarkGnalloyHTTPS2-8 320000 16554 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=935340 p95LatencyNs=2400506 p99LatencyNs=4399558 p999LatencyNs=7682315 maxLatencyNs=12230499 rssBytes=33947648 heapAllocBytes=13048944 heapSysBytes=23068672 heapObjects=135283 gcCount=148 gcPauseNs=58723762 goroutines=130 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=5.672141844s throughput=56416.08 ops/s
BenchmarkGnalloyHTTPS2-8 320000 17725 ns/op

framework=gnalloy protocol=https2 backend=epoll http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=900320 p95LatencyNs=1696610 p99LatencyNs=2418169 p999LatencyNs=4963584 maxLatencyNs=7211077 rssBytes=33804288 heapAllocBytes=7619256 heapSysBytes=22937600 heapObjects=23236 gcCount=151 gcPauseNs=28308211 goroutines=128 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=4.8150972s throughput=66457.64 ops/s
BenchmarkGnalloyHTTPS2-8 320000 15047 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 6.394894099s | 1 | 1 |  |
| 2 | 0 | 6.518064946s | 1 | 1 |  |
| 3 | 0 | 5.655383385s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 60407.78 | 915754 | 3665159 | 33906688 | 148 | 5.297330579s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 56416.08 | 935340 | 4399558 | 33947648 | 148 | 5.672141844s |
| gnalloy | https2 | h2 | epoll | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 66457.64 | 900320 | 2418169 | 33804288 | 151 | 4.8150972s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 16554 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 17725 | 0 | 0 |
| BenchmarkGnalloyHTTPS2-8 | 320000 | 15047 | 0 | 0 |

### netty epoll https2 tls13 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | https2 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 1m3.900890956s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol https2 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn h2 --tls-version 1.3
```

Output:

```text
framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=947859 p95LatencyNs=4617849 p99LatencyNs=8981866 p999LatencyNs=15602226 maxLatencyNs=16974174 rssBytes=587468800 heapAllocBytes=157651504 heapSysBytes=383778816 heapObjects=0 gcCount=20 gcPauseNs=80000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT7.640807118S throughput=41880.39 ops/s
BenchmarkNettyHTTPS2-8 320000 23878 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1086977 p95LatencyNs=4450834 p99LatencyNs=7169950 p999LatencyNs=13605606 maxLatencyNs=20175271 rssBytes=513060864 heapAllocBytes=78294576 heapSysBytes=318767104 heapObjects=0 gcCount=24 gcPauseNs=70000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT8.389439893S throughput=38143.19 ops/s
BenchmarkNettyHTTPS2-8 320000 26217 ns/op

framework=netty protocol=https2 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h2 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=976153 p95LatencyNs=5203200 p99LatencyNs=9003949 p999LatencyNs=14503770 maxLatencyNs=30724755 rssBytes=507047936 heapAllocBytes=111581064 heapSysBytes=318767104 heapObjects=0 gcCount=24 gcPauseNs=70000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT8.789172746S throughput=36408.43 ops/s
BenchmarkNettyHTTPS2-8 320000 27466 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 15.57334845s | 1 | 1 |  |
| 2 | 0 | 15.649579879s | 1 | 1 |  |
| 3 | 0 | 16.322204255s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | https2 | h2 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 41880.39 | 947859 | 8981866 | 587468800 | 20 | 7.640807118s |
| netty | https2 | h2 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 38143.19 | 1086977 | 7169950 | 513060864 | 24 | 8.389439893s |
| netty | https2 | h2 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 36408.43 | 976153 | 9003949 | 507047936 | 24 | 8.789172746s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTPS2-8 | 320000 | 23878 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 26217 | 0 | 0 |
| BenchmarkNettyHTTPS2-8 | 320000 | 27466 | 0 | 0 |

### gnet https2 unsupported

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | https2 |
| backend | unsupported |
| payload | - |
| duration | 136ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
gnet benchmark harness does not provide HTTP/2 over TLS
```

### fasthttp https2 unsupported

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | https2 |
| backend | unsupported |
| payload | - |
| duration | 92ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
fasthttp does not provide an HTTP/2 server implementation
```

### netpoll https2 unsupported

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | https2 |
| backend | unsupported |
| payload | - |
| duration | 77ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
netpoll benchmark harness does not provide HTTP/2 over TLS
```

### gnalloy rfc9000 http3 tls13 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | http3 |
| backend | rfc9000 |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 49.148096105s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol http3 -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn h3 -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=http3 backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=2120793 p95LatencyNs=3838075 p99LatencyNs=4456117 p999LatencyNs=4806134 maxLatencyNs=5172729 rssBytes=1506222080 heapAllocBytes=921729896 heapSysBytes=1454145536 heapObjects=11141449 gcCount=45 gcPauseNs=10319832 goroutines=7 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=11.337125677s throughput=28225.85 ops/s
BenchmarkGnalloyHTTP3-8 320000 35429 ns/op

framework=gnalloy protocol=http3 backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1934202 p95LatencyNs=3606314 p99LatencyNs=4137028 p999LatencyNs=4569934 maxLatencyNs=4743132 rssBytes=1416114176 heapAllocBytes=698015808 heapSysBytes=1366163456 heapObjects=8225491 gcCount=47 gcPauseNs=8607167 goroutines=8 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=10.861063119s throughput=29463.05 ops/s
BenchmarkGnalloyHTTP3-8 320000 33941 ns/op

framework=gnalloy protocol=http3 backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1894616 p95LatencyNs=3368497 p99LatencyNs=3878332 p999LatencyNs=4253030 maxLatencyNs=4667141 rssBytes=1459339264 heapAllocBytes=1123903032 heapSysBytes=1408040960 heapObjects=13803242 gcCount=46 gcPauseNs=12573830 goroutines=14 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=10.633613046s throughput=30093.25 ops/s
BenchmarkGnalloyHTTP3-8 320000 33230 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 12.977379446s | 1 | 1 |  |
| 2 | 0 | 12.426690282s | 1 | 1 |  |
| 3 | 0 | 12.016452304s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 28225.85 | 2120793 | 4456117 | 1506222080 | 45 | 11.337125677s |
| gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 29463.05 | 1934202 | 4137028 | 1416114176 | 47 | 10.861063119s |
| gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 30093.25 | 1894616 | 3878332 | 1459339264 | 46 | 10.633613046s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTP3-8 | 320000 | 35429 | 0 | 0 |
| BenchmarkGnalloyHTTP3-8 | 320000 | 33941 | 0 | 0 |
| BenchmarkGnalloyHTTP3-8 | 320000 | 33230 | 0 | 0 |

### netty epoll http3 tls13 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | http3 |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 2m22.822613725s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol http3 --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn h3 --tls-version 1.3
```

Output:

```text
framework=netty protocol=http3 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=4527401 p95LatencyNs=5641797 p99LatencyNs=10140036 p999LatencyNs=37198922 maxLatencyNs=116135432 rssBytes=423370752 heapAllocBytes=78127320 heapSysBytes=264241152 heapObjects=0 gcCount=20 gcPauseNs=88000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT24.664824847S throughput=12973.94 ops/s
BenchmarkNettyHTTP3-8 320000 77078 ns/op

framework=netty protocol=http3 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=4490056 p95LatencyNs=5990815 p99LatencyNs=9844845 p999LatencyNs=36493647 maxLatencyNs=44805777 rssBytes=412467200 heapAllocBytes=141910568 heapSysBytes=264241152 heapObjects=0 gcCount=19 gcPauseNs=92000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT24.57192803S throughput=13022.99 ops/s
BenchmarkNettyHTTP3-8 320000 76787 ns/op

framework=netty protocol=http3 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=4538072 p95LatencyNs=6370788 p99LatencyNs=10272673 p999LatencyNs=38079540 maxLatencyNs=41397841 rssBytes=412934144 heapAllocBytes=79638888 heapSysBytes=264241152 heapObjects=0 gcCount=20 gcPauseNs=77000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT25.061816195S throughput=12768.43 ops/s
BenchmarkNettyHTTP3-8 320000 78318 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 35.695477724s | 1 | 1 |  |
| 2 | 0 | 35.75225935s | 1 | 1 |  |
| 3 | 0 | 35.354513861s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | http3 | h3 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 12973.94 | 4527401 | 10140036 | 423370752 | 20 | 24.664824847s |
| netty | http3 | h3 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 13022.99 | 4490056 | 9844845 | 412467200 | 19 | 24.57192803s |
| netty | http3 | h3 | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 12768.43 | 4538072 | 10272673 | 412934144 | 20 | 25.061816195s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTP3-8 | 320000 | 77078 | 0 | 0 |
| BenchmarkNettyHTTP3-8 | 320000 | 76787 | 0 | 0 |
| BenchmarkNettyHTTP3-8 | 320000 | 78318 | 0 | 0 |

### gnalloy rfc9000 http3 tls13 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | http3 |
| backend | rfc9000 |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 48.066737052s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol http3 -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn h3 -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=http3 backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1899542 p95LatencyNs=3780456 p99LatencyNs=4500126 p999LatencyNs=5094080 maxLatencyNs=5358876 rssBytes=1343291392 heapAllocBytes=1147303864 heapSysBytes=1294827520 heapObjects=14122862 gcCount=45 gcPauseNs=23466703 goroutines=6 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=11.149341652s throughput=28701.25 ops/s
BenchmarkGnalloyHTTP3-8 320000 34842 ns/op

framework=gnalloy protocol=http3 backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1838935 p95LatencyNs=3576816 p99LatencyNs=4291853 p999LatencyNs=4890155 maxLatencyNs=5251870 rssBytes=1432805376 heapAllocBytes=1086341048 heapSysBytes=1382875136 heapObjects=13319799 gcCount=45 gcPauseNs=11987856 goroutines=14 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=10.37984988s throughput=30828.96 ops/s
BenchmarkGnalloyHTTP3-8 320000 32437 ns/op

framework=gnalloy protocol=http3 backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1856469 p95LatencyNs=3410563 p99LatencyNs=3848737 p999LatencyNs=4232269 maxLatencyNs=4674097 rssBytes=1471246336 heapAllocBytes=1147920296 heapSysBytes=1420591104 heapObjects=14133853 gcCount=44 gcPauseNs=8614702 goroutines=14 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=10.391318812s throughput=30794.94 ops/s
BenchmarkGnalloyHTTP3-8 320000 32473 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 12.779685957s | 1 | 1 |  |
| 2 | 0 | 11.795365333s | 1 | 1 |  |
| 3 | 0 | 11.750895695s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 28701.25 | 1899542 | 4500126 | 1343291392 | 45 | 11.149341652s |
| gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 30828.96 | 1838935 | 4291853 | 1432805376 | 45 | 10.37984988s |
| gnalloy | http3 | h3 | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 30794.94 | 1856469 | 3848737 | 1471246336 | 44 | 10.391318812s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyHTTP3-8 | 320000 | 34842 | 0 | 0 |
| BenchmarkGnalloyHTTP3-8 | 320000 | 32437 | 0 | 0 |
| BenchmarkGnalloyHTTP3-8 | 320000 | 32473 | 0 | 0 |

### netty epoll http3 tls13 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | http3 |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 2m25.158982391s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol http3 --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn h3 --tls-version 1.3
```

Output:

```text
framework=netty protocol=http3 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=4552072 p95LatencyNs=5726988 p99LatencyNs=9913636 p999LatencyNs=36901946 maxLatencyNs=97751785 rssBytes=421994496 heapAllocBytes=127955552 heapSysBytes=264241152 heapObjects=0 gcCount=23 gcPauseNs=69000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT24.707941247S throughput=12951.30 ops/s
BenchmarkNettyHTTP3-8 320000 77212 ns/op

framework=netty protocol=http3 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=4574474 p95LatencyNs=6172424 p99LatencyNs=9637870 p999LatencyNs=36386458 maxLatencyNs=41612037 rssBytes=425512960 heapAllocBytes=128345232 heapSysBytes=264241152 heapObjects=0 gcCount=23 gcPauseNs=51000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT25.007919644S throughput=12795.95 ops/s
BenchmarkNettyHTTP3-8 320000 78150 ns/op

framework=netty protocol=http3 backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=h3 latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=4604703 p95LatencyNs=6143361 p99LatencyNs=9605414 p999LatencyNs=37048108 maxLatencyNs=43780937 rssBytes=401076224 heapAllocBytes=106499808 heapSysBytes=264241152 heapObjects=0 gcCount=23 gcPauseNs=51000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT25.18779342S throughput=12704.57 ops/s
BenchmarkNettyHTTP3-8 320000 78712 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 35.998674155s | 1 | 1 |  |
| 2 | 0 | 36.122277863s | 1 | 1 |  |
| 3 | 0 | 36.146975148s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | http3 | h3 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 12951.3 | 4552072 | 9913636 | 421994496 | 23 | 24.707941247s |
| netty | http3 | h3 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 12795.95 | 4574474 | 9637870 | 425512960 | 23 | 25.007919644s |
| netty | http3 | h3 | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 12704.57 | 4604703 | 9605414 | 401076224 | 23 | 25.18779342s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyHTTP3-8 | 320000 | 77212 | 0 | 0 |
| BenchmarkNettyHTTP3-8 | 320000 | 78150 | 0 | 0 |
| BenchmarkNettyHTTP3-8 | 320000 | 78712 | 0 | 0 |

### gnet http3 unsupported

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | http3 |
| backend | unsupported |
| payload | - |
| duration | 164ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
gnet does not provide an HTTP/3 over QUIC codec in this benchmark harness
```

### fasthttp http3 unsupported

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | http3 |
| backend | unsupported |
| payload | - |
| duration | 96ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
fasthttp is HTTP/1 oriented and does not provide HTTP/3
```

### netpoll http3 unsupported

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | http3 |
| backend | unsupported |
| payload | - |
| duration | 84ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
netpoll benchmark harness does not provide HTTP/3 over QUIC
```

### gnalloy rfc9000 quic-stream tls13 128B

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | quic-stream |
| backend | rfc9000 |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 31.749631365s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol quic-stream -payload 128 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn gnalloy-quic -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=quic-stream backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1279754 p95LatencyNs=3067349 p99LatencyNs=5114955 p999LatencyNs=8544302 maxLatencyNs=11123402 rssBytes=37687296 heapAllocBytes=12598168 heapSysBytes=32309248 heapObjects=109376 gcCount=232 gcPauseNs=71262007 goroutines=6 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=7.413946059s throughput=43161.90 ops/s
BenchmarkGnalloyQUICStream-8 320000 23169 ns/op

framework=gnalloy protocol=quic-stream backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1234367 p95LatencyNs=2826492 p99LatencyNs=4756887 p999LatencyNs=6871461 maxLatencyNs=9627917 rssBytes=39854080 heapAllocBytes=11333728 heapSysBytes=32374784 heapObjects=87953 gcCount=235 gcPauseNs=76429896 goroutines=6 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=7.068079905s throughput=45273.96 ops/s
BenchmarkGnalloyQUICStream-8 320000 22088 ns/op

framework=gnalloy protocol=quic-stream backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1166709 p95LatencyNs=2709175 p99LatencyNs=4027335 p999LatencyNs=5728587 maxLatencyNs=6989243 rssBytes=37736448 heapAllocBytes=9608952 heapSysBytes=32342016 heapObjects=62659 gcCount=235 gcPauseNs=62682702 goroutines=6 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=6.732191305s throughput=47532.81 ops/s
BenchmarkGnalloyQUICStream-8 320000 21038 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 8.293187319s | 1 | 1 |  |
| 2 | 0 | 7.933891893s | 1 | 1 |  |
| 3 | 0 | 7.516413765s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 43161.9 | 1279754 | 5114955 | 37687296 | 232 | 7.413946059s |
| gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 45273.96 | 1234367 | 4756887 | 39854080 | 235 | 7.068079905s |
| gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 128 | 64 | 5000 | 320000 | 0 | 47532.81 | 1166709 | 4027335 | 37736448 | 235 | 6.732191305s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyQUICStream-8 | 320000 | 23169 | 0 | 0 |
| BenchmarkGnalloyQUICStream-8 | 320000 | 22088 | 0 | 0 |
| BenchmarkGnalloyQUICStream-8 | 320000 | 21038 | 0 | 0 |

### netty epoll quic-stream tls13 128B

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | quic-stream |
| backend | epoll |
| payload | 128B |
| warmup | 1 |
| repeat | 3 |
| duration | 38m45.512353878s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol quic-stream --backend epoll --payload 128 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn gnalloy-quic --tls-version 1.3
```

Output:

```text
framework=netty protocol=quic-stream backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=105610940 p95LatencyNs=202143837 p99LatencyNs=385753952 p999LatencyNs=637052047 maxLatencyNs=1059402896 rssBytes=1238790144 heapAllocBytes=376389640 heapSysBytes=476053504 heapObjects=0 gcCount=31 gcPauseNs=259000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT9M28.510315805S throughput=562.87 ops/s
BenchmarkNettyQUICStream-8 320000 1776595 ns/op

framework=netty protocol=quic-stream backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=105754949 p95LatencyNs=198960330 p99LatencyNs=380746959 p999LatencyNs=512986739 maxLatencyNs=1148231459 rssBytes=1247531008 heapAllocBytes=354096400 heapSysBytes=480247808 heapObjects=0 gcCount=31 gcPauseNs=267000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT9M30.0255379S throughput=561.38 ops/s
BenchmarkNettyQUICStream-8 320000 1781330 ns/op

framework=netty protocol=quic-stream backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=106240830 p95LatencyNs=200958300 p99LatencyNs=345269871 p999LatencyNs=446602911 maxLatencyNs=511696040 rssBytes=1244352512 heapAllocBytes=343750024 heapSysBytes=482344960 heapObjects=0 gcCount=31 gcPauseNs=261000000 goroutines=0 payload=128 connections=64 messages=5000 total=320000 errors=0 elapsed=PT9M32.033336759S throughput=559.41 ops/s
BenchmarkNettyQUICStream-8 320000 1787604 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 9m42.193283291s | 1 | 1 |  |
| 2 | 0 | 9m43.552319774s | 1 | 1 |  |
| 3 | 0 | 9m45.963948565s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 562.87 | 105610940 | 385753952 | 1238790144 | 31 | 9m28.510315805s |
| netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 561.38 | 105754949 | 380746959 | 1247531008 | 31 | 9m30.0255379s |
| netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 128 | 64 | 5000 | 320000 | 0 | 559.41 | 106240830 | 345269871 | 1244352512 | 31 | 9m32.033336759s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyQUICStream-8 | 320000 | 1776595 | 0 | 0 |
| BenchmarkNettyQUICStream-8 | 320000 | 1781330 | 0 | 0 |
| BenchmarkNettyQUICStream-8 | 320000 | 1787604 | 0 | 0 |

### gnalloy rfc9000 quic-stream tls13 1KiB

| Field | Value |
| --- | --- |
| framework | gnalloy |
| protocol | quic-stream |
| backend | rfc9000 |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 32.053782137s |
| exitCode | 0 |
| skipped | false |

Command:

```text
./benchmarks/external/bin/gnalloy-bench -protocol quic-stream -payload 1024 -connections 64 -messages 5000 -latency-sample-rate 64 -warmup-messages 500 -timeout 15m -alpn gnalloy-quic -tls-version 1.3
```

Output:

```text
framework=gnalloy protocol=quic-stream backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1182089 p95LatencyNs=2646358 p99LatencyNs=4166433 p999LatencyNs=7431565 maxLatencyNs=11882693 rssBytes=40570880 heapAllocBytes=12252832 heapSysBytes=32243712 heapObjects=93692 gcCount=281 gcPauseNs=69280389 goroutines=6 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=6.766262135s throughput=47293.47 ops/s
BenchmarkGnalloyQUICStream-8 320000 21145 ns/op

framework=gnalloy protocol=quic-stream backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1195534 p95LatencyNs=2614530 p99LatencyNs=4004240 p999LatencyNs=5887136 maxLatencyNs=12311039 rssBytes=38965248 heapAllocBytes=10696456 heapSysBytes=28082176 heapObjects=65219 gcCount=278 gcPauseNs=72114789 goroutines=6 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=6.786800755s throughput=47150.35 ops/s
BenchmarkGnalloyQUICStream-8 320000 21209 ns/op

framework=gnalloy protocol=quic-stream backend=rfc9000 http1Mode=codec boss=1 workers=4 readBufferSize=4096 reuseport=false mmap=false mmapBlockSize=4096 mmapBlocks=4096 iouringFixedBuffers=false iouringMultishotAccept=false iouringSQPoll=false tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=1244600 p95LatencyNs=2957165 p99LatencyNs=4482656 p999LatencyNs=7445889 maxLatencyNs=13051579 rssBytes=39247872 heapAllocBytes=11094896 heapSysBytes=32210944 heapObjects=75739 gcCount=277 gcPauseNs=93951167 goroutines=6 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=7.101220499s throughput=45062.68 ops/s
BenchmarkGnalloyQUICStream-8 320000 22191 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 7.668845303s | 1 | 1 |  |
| 2 | 0 | 7.776724955s | 1 | 1 |  |
| 3 | 0 | 8.079239222s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 47293.47 | 1182089 | 4166433 | 40570880 | 281 | 6.766262135s |
| gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 47150.35 | 1195534 | 4004240 | 38965248 | 278 | 6.786800755s |
| gnalloy | quic-stream | gnalloy-quic | rfc9000 | boss=1 workers=4 readBuffer=4096 mmapBlock=4096 mmapBlocks=4096 latencySampleRate=64 | 1024 | 64 | 5000 | 320000 | 0 | 45062.68 | 1244600 | 4482656 | 39247872 | 277 | 7.101220499s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkGnalloyQUICStream-8 | 320000 | 21145 | 0 | 0 |
| BenchmarkGnalloyQUICStream-8 | 320000 | 21209 | 0 | 0 |
| BenchmarkGnalloyQUICStream-8 | 320000 | 22191 | 0 | 0 |

### netty epoll quic-stream tls13 1KiB

| Field | Value |
| --- | --- |
| framework | netty |
| protocol | quic-stream |
| backend | epoll |
| payload | 1KiB |
| warmup | 1 |
| repeat | 3 |
| duration | 32m41.602811309s |
| exitCode | 0 |
| skipped | false |

Command:

```text
java -jar ./benchmarks/external/bin/netty-bench.jar --protocol quic-stream --backend epoll --payload 1024 --connections 64 --messages 5000 --latency-sample-rate 64 --warmup-messages 500 --timeout 15m --alpn gnalloy-quic --tls-version 1.3
```

Output:

```text
framework=netty protocol=quic-stream backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=72604009 p95LatencyNs=248060249 p99LatencyNs=510589424 p999LatencyNs=1166554599 maxLatencyNs=1689819850 rssBytes=1153892352 heapAllocBytes=384142216 heapSysBytes=429916160 heapObjects=0 gcCount=30 gcPauseNs=223000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT7M53.504538557S throughput=675.81 ops/s
BenchmarkNettyQUICStream-8 320000 1479702 ns/op

framework=netty protocol=quic-stream backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=73157294 p95LatencyNs=260252242 p99LatencyNs=545051977 p999LatencyNs=1613607811 maxLatencyNs=3407826445 rssBytes=1157861376 heapAllocBytes=384322888 heapSysBytes=434110464 heapObjects=0 gcCount=30 gcPauseNs=241000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT8M1.897604131S throughput=664.04 ops/s
BenchmarkNettyQUICStream-8 320000 1505930 ns/op

framework=netty protocol=quic-stream backend=epoll eventLoops=8 tlsVersion=1.3 cipherSuites= negotiatedProtocol=gnalloy-quic latencySampleRate=64 warmupMessages=500 latencySamples=5056 p50LatencyNs=73939979 p95LatencyNs=246787733 p99LatencyNs=545080776 p999LatencyNs=1121477449 maxLatencyNs=1554828934 rssBytes=1153978368 heapAllocBytes=284904584 heapSysBytes=427819008 heapObjects=0 gcCount=31 gcPauseNs=243000000 goroutines=0 payload=1024 connections=64 messages=5000 total=320000 errors=0 elapsed=PT8M0.218558545S throughput=666.36 ops/s
BenchmarkNettyQUICStream-8 320000 1500683 ns/op

```

Samples:

| Index | Exit | Duration | Stats | Metrics | Error |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 8m5.983322131s | 1 | 1 |  |
| 2 | 0 | 8m15.457526814s | 1 | 1 |  |
| 3 | 0 | 8m12.69243541s | 1 | 1 |  |

Stats:

| Framework | Protocol | ALPN | Backend | Loops | Payload | Connections | Messages | Total | Errors | Throughput ops/s | P50 latency ns | P99 latency ns | RSS bytes | GC count | Elapsed |
| --- | --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 675.81 | 72604009 | 510589424 | 1153892352 | 30 | 7m53.504538557s |
| netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 664.04 | 73157294 | 545051977 | 1157861376 | 30 | 8m1.897604131s |
| netty | quic-stream | gnalloy-quic | epoll | eventLoops=8 | 1024 | 64 | 5000 | 320000 | 0 | 666.36 | 73939979 | 545080776 | 1153978368 | 31 | 8m0.218558545s |

Metrics:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| BenchmarkNettyQUICStream-8 | 320000 | 1479702 | 0 | 0 |
| BenchmarkNettyQUICStream-8 | 320000 | 1505930 | 0 | 0 |
| BenchmarkNettyQUICStream-8 | 320000 | 1500683 | 0 | 0 |

### gnet quic stream unsupported

| Field | Value |
| --- | --- |
| framework | gnet |
| protocol | quic-stream |
| backend | unsupported |
| payload | - |
| duration | 113ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
gnet does not provide a QUIC transport
```

### fasthttp quic stream unsupported

| Field | Value |
| --- | --- |
| framework | fasthttp |
| protocol | quic-stream |
| backend | unsupported |
| payload | - |
| duration | 64ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
fasthttp does not provide a QUIC transport
```

### netpoll quic stream unsupported

| Field | Value |
| --- | --- |
| framework | netpoll |
| protocol | quic-stream |
| backend | unsupported |
| payload | - |
| duration | 61ns |
| exitCode | 0 |
| skipped | true |

Command:

```text

```

Output:

```text
netpoll does not provide a QUIC transport
```
