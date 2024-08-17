/*
    https://developer.mozilla.org/en-US/docs/Web/API/Document/querySelector
    https://developer.mozilla.org/en-US/docs/Web/API/Element/keydown_event
    https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/key
    https://developer.mozilla.org/en-US/docs/Web/API/Element/click_event
    https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio
    https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement
    https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API/Using_the_Web_Storage_API
*/

let tracks = [];
let trackNumber = 0;

function pauseOrPlay() {
    const audioPlayer = document.querySelector("#player");
    const playButton = document.querySelector("#play");

    if (audioPlayer.paused === true) {
        playButton.textContent = "Pause";
        audioPlayer.play();
    } else {
        playButton.textContent = "Play";
        audioPlayer.pause();
    }
}

function changeTrack(n, oldN) {
    const track = tracks[n];

    const audioPlayer = document.querySelector("#player");
    const audioPlayerSource = document.querySelector("#playersource");
    const nameLabel = document.querySelector("#name");

    nameLabel.textContent = track["name"];
    audioPlayerSource.src = track["source"];
    audioPlayerSource.type = track["mimeType"];

    const oldTrackElement = document.querySelector("#track" + oldN)
    const newTrackElement = document.querySelector("#track" + n)

    oldTrackElement.className = "clickable-track"
    newTrackElement.className = "active-track"

    const paused = audioPlayer.paused;
    audioPlayer.load();
    if (paused !== true) {
        audioPlayer.play();
    }
}

function previousTrack() {
    let oldTrackNumber = trackNumber;

    if (trackNumber === 0) {
        trackNumber = tracks.length - 1;
    } else {
        trackNumber--;
    }
    changeTrack(trackNumber, oldTrackNumber);
}

function nextTrack() {
    let oldTrackNumber = trackNumber;

    trackNumber++;
    if (trackNumber >= tracks.length) {
        trackNumber = 0;
    }
    changeTrack(trackNumber, oldTrackNumber);
}

function changePosition(n) {
    const audioPlayer = document.querySelector("#player");

    audioPlayer.currentTime += n;
}

// initAudioPlayer is called from the HTML generated from playlistsbyid.tmpl.html
// eslint-disable-next-line no-unused-vars
function initAudioPlayer(availableTracks, n) {
    // Initial state
    tracks = availableTracks;
    trackNumber = n;
    changeTrack(0, n);

    // Set up events
    const previousButton = document.querySelector("#previous");
    const playButton = document.querySelector("#play");
    const nextButton = document.querySelector("#next");
    const audioPlayer = document.querySelector("#player");

    // Load/save volume level in local storage (where available).
    loadAudio(audioPlayer);
    audioPlayer.addEventListener("volumechange", () => {
        saveAudio(audioPlayer);
    });

    // Play/pause via audio player or button
    playButton.addEventListener("click", () => {
        pauseOrPlay();
    });
    audioPlayer.addEventListener("pause", () => {
        playButton.textContent = "Play";
    });
    audioPlayer.addEventListener("play", () => {
        playButton.textContent = "Pause";
    });

    // Buttons to move back/forward
    previousButton.addEventListener("click", () => {
        previousTrack();
    });

    nextButton.addEventListener("click", () => {
        nextTrack();
    });

    // Automatically move onto next track when one finishes.
    // Stop if this was the last track.
    audioPlayer.addEventListener("ended", () => {
        if (trackNumber + 1 >= tracks.length) {
            return;
        }

        nextTrack();
        // When it ends, it seems to go into the paused state.
        // But we know it was playing, otherwise the track would not have ended!
        audioPlayer.play();
    });

    // Jump to a track if its name is clicked.
    for (let i = 0; i < tracks.length; ++i) {
        const trackElement = document.querySelector("#track" + i);
        trackElement.addEventListener("click", () => {
            let oldTrackNumber = trackNumber;
            trackNumber = i;
            changeTrack(trackNumber, oldTrackNumber);
        });
    }

    // Hotkeys
    document.addEventListener("keydown", (event) => {
        if (event.isComposing) {
            return;
        }

        // If the play button is already the target of events,
        // ignore this event. Let the button process its keypresses
        // as normal. This avoids "bouncing" where space would trigger
        // both a click and keydown event, resulting in pause-unpause
        // (or unpause-pause).
        if (event.key === " " && event.target.id !== "play") { // space
            pauseOrPlay();
            return;
        }
        if (event.key === ",") { // previous
            previousTrack();
            return;
        }
        if (event.key === ".") { // next
            nextTrack();
            return;
        }
        if (event.key === "<") { // rewind
            changePosition(-15.0);
            return;
        }
        if (event.key === ">") { // fast-forward
            changePosition(15.0);
            return;
        }
    });
}

// Based on storageAvailable from
// https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API/Using_the_Web_Storage_API
function localStorageAvailable() {
    let storage;
    try {
        storage = window["localStorage"];
        const x = "__storage_test__";
        storage.setItem(x, x);
        storage.removeItem(x);
        return true;
    } catch (e) {
        return (
            e instanceof DOMException &&
            e.name === "QuotaExceededError" &&
            // acknowledge QuotaExceededError only if there's something already stored
            storage &&
            storage.length !== 0
        );
    }
}

// TODO: button to clear local storage
function loadAudio(audioPlayer) {
    let audio = {};

    if (localStorageAvailable()) {
        const val = localStorage.getItem("audio");
        if (val !== null) {
            audio = JSON.parse(val);
        }
    }

    if (audio.volume !== null && audio.volume !== undefined) {
        audioPlayer.volume = audio.volume;
    }
    if (audio.muted !== null && audio.muted !== undefined) {
        audioPlayer.muted = audio.muted;
    }
}

function saveAudio(audioPlayer) {
    const audio = {
        volume: audioPlayer.volume,
        muted: audioPlayer.muted,
    };
    const val = JSON.stringify(audio);

    if (localStorageAvailable()) {
        localStorage.setItem("audio", val);
    }
}