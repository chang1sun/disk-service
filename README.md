# EasyDisk
*The back-end repository for EasyDisk*

## Repo Structure
```
- analysis // provide analysis support according to user's behavior
    - resource   // proto file, managed by buf
    - stub       // pb deps, generate by protoc
    - common     // common constants
    - util       // utitity functions
    - interfaces // grpc interfaces
    - service    // service logic implementaion
    - repo       // interact with data
    - rpc        // a series of function make rpc calls
- auth     // provide authentication support
    - ...
- file     // provide basic file operation support
    - ...
- gateway  // RESTful api gateway, forwards http request from cli
    - ...
```