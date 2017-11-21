Thank you for contributing to meli.                    
Every contribution to meli is important to us.                                   

Contributor offers to license certain software (a “Contribution” or multiple
“Contributions”) to meli, and meli agrees to accept said Contributions,
under the terms of the MIT License.
Contributor understands and agrees that meli shall have the irrevocable and perpetual right to make
and distribute copies of any Contribution, as well as to create and distribute collective works and
derivative works of any Contribution, under the MIT License.

## To contribute:            

- fork this repo.
- make the changes you want on your fork.
- your changes should have backward compatibility in mind unless it is impossible to do so.
- add your name and contact(optional) to CONTRIBUTORS.md
- add tests and benchmarks
- format your code using gofmt:                                          
- run tests(with race flag) and make sure everything is passing:
```shell
 go test -race -cover -v ./...
```
- run benchmarks and make sure that they havent regressed. If you have introduced any regressions, fix them unless it is impossible to do so:
```shell
go test -race -run=XXXX -bench=. ./...
```
- open a pull request on this repo.          
          
NB: I make no commitment of accepting your pull requests.                 
