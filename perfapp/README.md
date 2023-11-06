# Performance application

## Available workloads

- Health: Returns 200. Does not need postgres connectivity. Available at `/health`
- Ready: Inserts a timestamp record in the ts table. Available at `/ready`
- Euler aproximation: Computes an Euler number approximation and writes compute time in the euler table. Available at `/euler`
