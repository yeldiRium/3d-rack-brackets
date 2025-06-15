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

	tmux split-window -t "${session}:0.0" -h -l "70%" -b

	tmux split-window -t "${session}:0.1" -v -l "20%" -b
	tmux split-window -t "${session}:0.2" -v -l "10%" -b
	tmux split-window -t "${session}:0.3" -v -l "30%" -b

	sleep 1

	tmux send-keys -t "${session}:0.0" " devbox shell" "C-m" "C-l"
	tmux send-keys -t "${session}:0.1" " devbox shell" "C-m" "C-l"
	tmux send-keys -t "${session}:0.2" " devbox shell" "C-m" "C-l"
	tmux send-keys -t "${session}:0.3" " devbox shell" "C-m" "C-l"
	tmux send-keys -t "${session}:0.4" " devbox shell" "C-m" "C-l"

	tmux send-keys -t "${session}:0.0" " nvim" "C-m" "M-1"
	tmux send-keys -t "${session}:0.1" " git bug termui" "C-m"
	tmux send-keys -t "${session}:0.2" " devbox run open" "C-m"
	tmux send-keys -t "${session}:0.3" " devbox run watch" "C-m"

	tmux select-pane -t "${session}:0.0"
fi

attach

