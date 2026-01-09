# Grip shell integration for Zsh
# Added via: grip shell >> ~/.zshrc

# Change directory to selected project
gpcd() {
    local result
    result=$(grip cd "$@")
    if [[ $? -eq 0 ]] && [[ -n "$result" ]]; then
        cd "$result"
    fi
}

# Switch project and branch, then cd
gpsw() {
    local result
    result=$(grip switch "$@" 2>&1 | tail -1)
    if [[ $? -eq 0 ]] && [[ -d "$result" ]]; then
        cd "$result"
    fi
}

# Switch branch in current repo
gpb() {
    grip branches "$@"
}

# Aliases
alias gps='grip status'
alias gpi='grip init'
