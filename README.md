# race-condition

In the output below, the first time we check for the container ports we get `56534`, `32877`, and `56536`.

After sleeping 3 seconds, we check again and get `56534`, `56535`, and `56536`.

Output:

```
2023/09/12 14:41:34 (Before sleeping)
2023/09/12 14:41:34 Total ports: 2
2023/09/12 14:41:34 Ctr ports: [{ 80 0 tcp} {0.0.0.0 8080 56534 tcp}]
2023/09/12 14:41:34 Host port: 0
2023/09/12 14:41:34 (Before sleeping)
2023/09/12 14:41:34 Total ports: 2
2023/09/12 14:41:34 Ctr ports: [{ 80 0 tcp} {127.0.0.1 8080 32877 tcp}]
2023/09/12 14:41:34 Host port: 0
2023/09/12 14:41:34 (Before sleeping)
2023/09/12 14:41:34 Total ports: 2
2023/09/12 14:41:34 Ctr ports: [{ 80 0 tcp} {0.0.0.0 8080 56536 tcp}]
2023/09/12 14:41:34 Host port: 0
2023/09/12 14:41:37 (After sleeping)
2023/09/12 14:41:37 Total ports: 2
2023/09/12 14:41:37 Ctr ports: [{ 80 0 tcp} {0.0.0.0 8080 56534 tcp}]
2023/09/12 14:41:37 Host port: 0
2023/09/12 14:41:37 (After sleeping)
2023/09/12 14:41:37 Total ports: 2
2023/09/12 14:41:37 Ctr ports: [{0.0.0.0 8080 56535 tcp} { 80 0 tcp}]
2023/09/12 14:41:37 Host port: 56535
2023/09/12 14:41:37 (After sleeping)
2023/09/12 14:41:37 Total ports: 2
2023/09/12 14:41:37 Ctr ports: [{ 80 0 tcp} {0.0.0.0 8080 56536 tcp}]
2023/09/12 14:41:37 Host port: 0
```
