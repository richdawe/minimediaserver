 * Better favicon for showing in browser title bar

 * Structured logging using the standard log package or something from https://blog.logrocket.com/5-structured-logging-packages-for-go/ ?

 * Listen to media hotkeys in browser (and optionally allow this to be disabled?)
 * Stop playing music when the computer is being suspended/hibernated

 * Document hotkeys somewhere

 * Bug: Still can't skip on long tracks when using remote minimediaserver (not local one)
   * This seems to be a Firefox HTML audio player thing, rather than a golang HTTP server thing.
   * Would it help to be able to support fetching data ranges on a track? Instead of streaming the response with echo .Stream()?
   * JS code could pre-fetch, to populate the browser's cache?
   * Completely override fetching the data in JS, and just give the player a buffer of data instead of a URL?

 * Tags:
   * FLAC (Vorbis comments) (DONE, needs test coverage)
   * Ogg (Vorbis comments) (DONE, needs test coverage)
   * ID3
   * ID3 v2
   * Use for artist + title instead of filenames in playlists
   * Generate playlists based off tags (if present) rather than file location - playlist per album

 * Save volume level across invocations of page? (cookies? local storage?)

 * Try out on mobile phone. Need media queries to adjust layout for smaller screens or readability?

 * Memory profiling
 * Coverage
 * Function documentation

 * Track data JSON blob - fetch that via API rather than including in generated HTML data
   * with OpenAPI schema and validation in golang code
 * Optimize track storage in media server (*Track instead of Track)
