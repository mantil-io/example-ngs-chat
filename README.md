# example-ngs-chat

Example of ngs and Mantil.

Using [example](https://github.com/nats-io/nats.ws/tree/master/examples) from
[nats.ws](https://github.com/nats-io/nats.ws) project. All html and JavaScript
are just slightly modified version of nats.ws
[chat](https://github.com/nats-io/nats.ws/blob/master/examples/chat.js) example.

Instead of posting message to the nats subject we are sending it to the backend.
There we have stored state of all messages. When client connects it gets list of
all messages from the backend. Backend is now place to add additional logic, for
example to moderate messages.

Backend currently holds chat state in memory so on each Lambda function restart
we will start with empty state. For testing purpose this works just fine. 

Both backend and client needs credentials to connect to ngs. Backend needs
credentials in `api/chat/tester.creds` file and client in `api/chat/tester.creds`.

[Here](https://synadia.com/ngs/signup) is described how to get required tools
(ngs, nsc) and how to create ngs free tier account.

Once you have tools (ngs, nsc) installed this will create credentials for tester
user and put them in the right place.

``` sh
nsc add operator -u synadia
nsc add account myaccount
nsc add user tester
nsc generate creds -a myaccount -n tester -o ./tester.creds

cp tester.creds public
cp tester.creds api/chat
```

After that `mantil deploy` and you are up.
