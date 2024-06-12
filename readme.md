### Easyflow-Backend rewritten in Golang
Not much documentation yet. Important for contribution:
```
sh setup.sh
```
Or 
```
pre-commit install
```
This ensures that you install the pre-commit hook that ensures
that the module is tidy after you installed librarys and lints the code.

### Conventions
- You added a new env variable? 
Add it to config.go so you can easily and safely access it later.
- You built a Middleware which should only work selectively on specific routes?
Create a routing group