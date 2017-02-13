# BACKER
```
____          _____ _  ________ _____  
|  _ \   /\   / ____| |/ /  ____|  __ \ 
| |_) | /  \ | |    | ' /| |__  | |__) |
|  _ < / /\ \| |    |  < |  __| |  _  / 
| |_) / ____ \ |____| . \| |____| | \ \ 
|____/_/    \_\_____|_|\_\______|_|  \_\
```

### Architecture

Current approach assumes existing application Server <--> Client with RPC connection to make possible management
connection between these daemons and data communications through sockets.
Server is able to invoke operation on client side using RPC. Basic operations at this moment:
- Check connection (Ping-Pong)
- Run Backup (Make archive on client side and send it by data connection)

These operations are under development, so some of them can be changed

In this approach the application has limitation which is mirroring data on client side.
Client daemon creates archive with given paths in temporary path and then send the archive by socket connection to Server.




### TODO List

- [x] Prepare basic client - server structure
- [x] Extract common function to one package
- [x] Change 'prints' to log at the begging, later make it more configurable
- [x] BACSRV: Create package to read configuration files
- [x] BACSRV: Demonize server (It is not needed because of systemd)
- [x] BACSRV: Read config file from parameter
- [ ] BACSRV: Create inside gRPC API -- in progress
- [ ] BACLNT: Create inside gRPC API -- in progress
- [ ] BACSRV: Create oudside REST API -- in progress
- [ ] README: Add simple architecture picture in ASCI
- [ ] BACSRV: Prepare description of data repository structure
- [ ] BACSRV: Prepare description of database structure
- [ ] BACSRV: Create package to be able to write metadata of backed up files
