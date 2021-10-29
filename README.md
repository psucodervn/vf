## Install
```sh
curl https://raw.githubusercontent.com/psucodervn/vf/master/scripts/install.sh | sh
``` 

## Add to bash
```sh
v () {
    local dir=$(vf $@) 
    if [ -d "$dir" ]
    then
        cd "$dir"
    fi
}
```
