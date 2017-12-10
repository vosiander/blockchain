# blockchain
A basic blockchain implementation in go

## How to start

The blockchain can be created and maintained via the JSON-REST API by default running on Port 8080.

### Show all blocks
    
    curl -X GET http://localhost:8080/blocks 
    
### Create a new block
    
    curl -X POST http://localhost:8080/blocks -H 'Content-Type: application/json' -d '{
     "data": "And there is some more data here. A lot more data!"
    }' 
   
## Resources
A list of online resources I found.

- https://github.com/gin-gonic/gin/blob/master/README.md  Gin Router for the http part
- https://github.com/lhartikk/naivechain Naive blockchain implementation
- https://github.com/conradoqg/naivecoin Naive coin implementation