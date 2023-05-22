 * Better favicon for showing in browser title bar

 * Structured logging using the standard log package or something from https://blog.logrocket.com/5-structured-logging-packages-for-go/ ?

 * Listen to media hotkeys in browser (and optionally allow this to be disabled?)
 * Stop playing music when the computer is being suspended/hibernated

 * Bug: Space doesn't seem pause/play correctly after changing track, when mouse is over media controls - need to be aware of mouseenter/mouseleave so space doesn't get double-processed?
 * Document hotkeys somewhere

 * Tags:
   * FLAC (Vorbis comments) (DONE, needs test coverage)
   * Ogg (Vorbis comments) (DONE, needs test coverage)
   * ID3
   * ID3 v2

 * Save volume level across invocations of page? (cookies? local storage?)

 * Coverage
 * Function documentation

 * Track data JSON blob - fetch that via API rather than including in generated HTML data
   * with OpenAPI schema and validation in golang code
 * Optimize track storage in media server (*Track instead of Track)