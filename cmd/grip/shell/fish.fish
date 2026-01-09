# Grip shell integration for Fish
# Added via: grip shell --fish >> ~/.config/fish/config.fish

# Change directory to selected project
function gpcd
    set result (grip cd $argv)
    if test $status -eq 0; and test -n "$result"
        cd $result
    end
end

# Switch project and branch, then cd
function gpsw
    set result (grip switch $argv 2>&1 | tail -1)
    if test $status -eq 0; and test -d "$result"
        cd $result
    end
end

# Switch branch in current repo
function gpb
    grip branches $argv
end

# Aliases
alias gps='grip status'
alias gpi='grip init'
