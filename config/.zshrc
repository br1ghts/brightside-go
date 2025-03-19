# ==============================
# 🚀 Brightside-Go Zsh Configuration
# ==============================

# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

# ==============================
# 🚀 Brightside-Go Zsh Configuration
# ==============================

# Set Zsh as the default shell
export SHELL=$(which zsh)

# Enable Powerlevel10k Theme
export ZSH="$HOME/.oh-my-zsh"
ZSH_THEME="powerlevel10k/powerlevel10k"

# Load Oh My Zsh
if [ -f "$ZSH/oh-my-zsh.sh" ]; then
    source "$ZSH/oh-my-zsh.sh"
else
    echo "⚠️ Warning: Oh My Zsh not found!"
fi

# Plugins
plugins=(git zsh-syntax-highlighting zsh-autosuggestions)

# Enable Powerlevel10k Instant Prompt
typeset -g POWERLEVEL9K_INSTANT_PROMPT=off
if [ -f ~/.p10k.zsh ] && [ -z "$P10K_SOURCED" ]; then
  export P10K_SOURCED=1
  source ~/.p10k.zsh
fi

# Custom Cyberpunk Welcome Message
function badass_welcome() {
    [[ $- == *i* ]] || return  # Only show in interactive shells

    emulate -L zsh
    local symbols=("█▒░" "▓█░" "░▒▓" "███" "▓▓▒" "▒░▓")
    for i in {1..10}; do
        print -Pn "%F{green}${symbols[RANDOM % ${#symbols[@]}]}%F{black} Initializing...\r"
        sleep 0.1
    done
    print ""  # Newline after animation

    sleep 0.2
    print -P "%F{cyan}[LOG] %F{green}Neural uplink established...%f"
    sleep 0.2
    print -P "%F{cyan}[LOG] %F{yellow}Decrypting access protocols...%f"
    sleep 0.2
    print -P "%F{cyan}[LOG] %F{magenta}Environment loaded successfully.%f"
    sleep 0.2

    print -P '%F{red}💀 SYSTEM ONLINE...%f'
    sleep 0.1
    print -P '%F{blue}🕶️ Welcome back to the Brightside, MrBrightside.%f'
}

# Show Welcome Message
[[ -n "$PS1" && -z "$WELCOME_RAN" ]] && { export WELCOME_RAN=1; badass_welcome; }

# History Settings
setopt HIST_IGNORE_ALL_DUPS SHARE_HISTORY INC_APPEND_HISTORY
HISTSIZE=10000
SAVEHIST=10000
HISTFILE=~/.zsh_history

# Auto-Correction & Navigation
setopt AUTO_CD CORRECT

# Aliases
alias c='clear'
alias ls='ls --color=auto'
alias ll='ls -lhF --color=auto'
alias la='ls -A --color=auto'
alias lla='ls -lha --color=auto'

# Git Aliases
alias gst='git status'
alias gaa='git add .'
alias gcm='git commit -m'
alias gco='git checkout'

# Quick Reload & Edits
alias zshrc='vim ~/.zshrc'
alias reload='source ~/.zshrc'

# Add Brightside-Go to Path
export PATH="/usr/local/bin:$PATH"