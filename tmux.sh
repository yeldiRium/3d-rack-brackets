#!/usr/bin/env bash
set -eu

session="ide-3d-rack-brackets"
attach() {
  if [ -n "${TMUX:-}" ]; then
    tmux switch-client -t "=${session}"
  else
    tmux attach-session -t "=${session}"
  fi
}

if ! tmux has-session -t "${session}" 2>/dev/null; then
  tmux new-session -d -s "${session}" -x "$(tput cols)" -y "$(tput lines)"

  tmux split-window -t "${session}:0.0" -h
  tmux resize-pane -t "${session}:0.0" -x "70%"

	tmux split-window -t "${session}:0.1" -v
	tmux split-window -t "${session}:0.2" -v

	tmux resize-pane -t "${session}:0.1" -y "10%"
	tmux resize-pane -t "${session}:0.2" -y "40%"
	tmux resize-pane -t "${session}:0.3" -y "50%"

  sleep 1

  tmux send-keys -t "${session}:0.0" "devbox shell" "C-m" "C-l"
  tmux send-keys -t "${session}:0.1" "devbox shell" "C-m" "C-l"
  tmux send-keys -t "${session}:0.2" "devbox shell" "C-m" "C-l"
  tmux send-keys -t "${session}:0.3" "devbox shell" "C-m" "C-l"

  tmux send-keys -t "${session}:0.0" "nvim" "C-m" "M-1"
	tmux send-keys -t "${session}:0.1" "devbox run open" "C-m"
	tmux send-keys -t "${session}:0.2" "devbox run watch" "C-m"

  tmux select-pane -t "${session}:0.0"
fi

attach

