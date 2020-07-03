# file-worker

This project only demonstrates my ability to write code.
 Usually I create one more layer for business logic of an application (between view layer and storage), but in this case it useless, 
 so I just passed the storage layer to view layer.
 
> Note: I didn't use any database/storage to make duplicated words filtering, I used just memory implementation, but added the abstraction (interface `UniqueMatcher`) to have the ability to do it.
> For this case I would use Redis sets with expiration or Aerospike.
 
### To Run test
Please use `make unit_test` command to do it.    

### To Run project
Please use `docker build --tag file-worker:0.1 . && docker run file-worker:0.1` command to do it.