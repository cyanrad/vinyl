

// adds seconds to the current time of the audio element
function getMovedAudioTime(currentTime: number, duration: number, seconds: number) {
    return Math.max(0, Math.min(currentTime + seconds, duration));
}

function getNewAudioTime(currentTime: number, duration: number, newTime: number) {
    return Math.max(0, Math.min(newTime, duration));
}

export { getMovedAudioTime, getNewAudioTime };