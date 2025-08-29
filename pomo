#!/usr/bin/env bash

STATE_FILE="/tmp/pomo_state"
PID_FILE="/tmp/pomo_pid"
HELP_MSG="Usage: pomo <study:break[:long_break]> | stats | stop

Commands:
  stats             display an overview of the study sessions.
  stop              stop the pomodoro script.

Examples:
  Stating a session with 25min study time, 5min of break time, and 10min long break time
  >>> pomo 25:5:10

  displaying session stats
  >>> pomo stats

  stopping the background timer
  >>> pomo stop
"
notify() {
    # Popup persistente com botão OK
    zenity --info --title="Pomodoro" --text="$1" --no-wrap --width=250 --height=100 2>/dev/null
    # echo "[Pomodoro] $1"
}

run_loop() {
    local study=$1
    local brk=$2
    local long_break=$3
    local sessions=0
    local total_study=0

    while true; do
        # --- Study session ---
        sessions=$((sessions + 1))
        total_study=$((total_study + study))

        echo "state=study" > "$STATE_FILE"
        echo "end=$(( $(date +%s) + study*60 ))" >> "$STATE_FILE"
        echo "sessions=$sessions" >> "$STATE_FILE"
        echo "total_study=$total_study" >> "$STATE_FILE"

        notify "📚 Session #$sessions — Study for $study minutes"
        sleep "${study}m" || exit 0

        # --- Break session ---
        if (( sessions % 4 == 0 )); then
            echo "state=long_break" > "$STATE_FILE"
            echo "end=$(( $(date +%s) + long_break*60 ))" >> "$STATE_FILE"
            echo "sessions=$sessions" >> "$STATE_FILE"
            echo "total_study=$total_study" >> "$STATE_FILE"

            notify "🌴 Long break time: $long_break minutes (after 4 sessions)"
            sleep "${long_break}m" || exit 0
        else
            echo "state=break" > "$STATE_FILE"
            echo "end=$(( $(date +%s) + brk*60 ))" >> "$STATE_FILE"
            echo "sessions=$sessions" >> "$STATE_FILE"
            echo "total_study=$total_study" >> "$STATE_FILE"

            notify "☕ Short break: $brk minutes"
            sleep "${brk}m" || exit 0
        fi
    done
}

start() {
    if [[ -f "$PID_FILE" ]] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
        echo "Pomodoro already running."
        exit 1
    fi

    IFS=":" read -r study brk long_break <<< "$1"
    long_break=${long_break:-15}   # se não informar, usa 15 min por padrão

    (run_loop "$study" "$brk" "$long_break") &   # libera shell
    echo $! > "$PID_FILE"
    echo "Pomodoro started: $study min study, $brk min break, $long_break min long break (every 4 sessions)"
}

stats() {
    if [[ ! -f "$STATE_FILE" ]]; then
        echo "No Pomodoro running."
        exit 0
    fi

    source "$STATE_FILE"
    local now=$(date +%s)
    local remaining=$((end - now))
    if (( remaining < 0 )); then
        echo "⚠️ Session ended but state not updated yet."
        exit 0
    fi

    local min=$((remaining / 60))
    local sec=$((remaining % 60))

    echo "------ Pomodoro Stats ------"
    echo "✅ Sessions finished: ${sessions:-0}"
    echo "⏳ Total study time: ${total_study:-0} min"
    case $state in
        study) echo "📚 Current: Studying" ;;
        break) echo "☕ Current: Short Break" ;;
        long_break) echo "🌴 Current: Long Break" ;;
    esac
    printf "⏱ Remaining: %02d:%02d\n" "$min" "$sec"
    echo "----------------------------"
}

stop() {
    if [[ -f "$PID_FILE" ]]; then
        kill "$(cat "$PID_FILE")" 2>/dev/null # && notify "Pomodoro stopped."
        rm -f "$PID_FILE" "$STATE_FILE"
        echo "Pomodoro stopped."
    else
        echo "No Pomodoro running."
    fi
}

case "$1" in
    *:*) start "$1" ;;
    stats) stats ;;
    stop) stop ;;
    *) printf "%s" "$HELP_MSG" ;; #"Usage: pomo <study:break[:long_break]> | stats | stop" ;;
esac

