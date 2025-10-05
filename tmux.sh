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

	tmux split-window -t "${session}:0.0" -v -l "70%" -b
	tmux split-window -t "${session}:0.0" -h -l "60%" -b

	tmux split-window -t "${session}:0.1" -v -l "70%" -b

	sleep 1

	tmux send-keys -t "${session}:0.0" " devenv shell zsh" "C-m" "C-l"
	tmux send-keys -t "${session}:0.1" " devenv shell zsh" "C-m" "C-l"
	tmux send-keys -t "${session}:0.2" " devenv shell zsh" "C-m" "C-l"
	tmux send-keys -t "${session}:0.3" " devenv shell zsh" "C-m" "C-l"

	tmux send-keys -t "${session}:0.0" " nvim" "C-m" "M-1"
	tmux send-keys -t "${session}:0.1" " devenv up watch" "C-m"
	tmux send-keys -t "${session}:0.2" " git bug termui" "C-m"

	tmux select-pane -t "${session}:0.0"
fi

attach
