# BACKER
```
____          _____ _  ________ _____  
|  _ \   /\   / ____| |/ /  ____|  __ \ 
| |_) | /  \ | |    | ' /| |__  | |__) |
|  _ < / /\ \| |    |  < |  __| |  _  / 
| |_) / ____ \ |____| . \| |____| | \ \ 
|____/_/    \_\_____|_|\_\______|_|  \_\
```

### Simple Architecture Overview

Current approach assumes existing application Server <--> Client with RPC connection to make possible management
connection between these daemons and data communications through sockets.
Server is able to invoke operation on client side using gRPC. An operations can be triggered from server ~~and client side~~. Either client or server have two listeners - data and management and can "talk" each other.

Main points for the first release:

- Client server basic authentication
- Basic encryption for the bftp protocol
- Backup specific paths from a client
- Restore data to a given client
- List available backups from client side

The points described above are ganeral. More information can be find in the documentation under specific part called release note.
 

### TODO List

- [x] Prepare basic client - server structure
- [x] Extract common function to one package
- [x] Change 'prints' to log at the begging, later make it more configurable
- [x] BACSRV: Create package to read configuration files
- [x] BACSRV: Demonize server (It is not needed because of systemd)
- [x] BACSRV: Read config file from parameter
- [x] README: Add simple architecture picture in ASCI
- [x] Make backup working!
- [x] BACLI: List backups
- [x] BACLI: Trigger backup
- [x] BACLI: Trigger restore with options (different remote path, one file restore)
- [x] Make restore working!
- [ ] Restore files metadata -- in progress
- [ ] Encryption
- [ ] Python testing framework -- in progress
- [ ] DOCS (manual doc update, create how to doc)
