 * Improve how it works in Chrome (MOSTLY DONE, needs more tests)
   * Can't seek in files
   * Chrome seems to download a bit and then halt
   * Support range queries?
   * golang net/http/fs has range parsing, but it's private, could borrow? https://cs.opensource.google/go/go/+/refs/tags/go1.21.3:src/net/http/fs.go;l=837
   * Also helpful: https://www.zeng.dev/post/2023-http-range-and-play-mp4-in-browser/
   * Maybe helpful? https://stackoverflow.com/questions/39961598/video-seeking-in-google-chrome-how-to-handle-continuous-partial-content-request

 * Add memory debugging endpoints
   * Stack dump to log too?

 * Basic video player support
   * Directory of videos => album

 * Media player tab doesn't give you an error when server is down; just hangs?

 * Command-line switches for help

 * Log stats during start-up
   * Time per storage backend to evaluate all files
   * Total number of tracks found, ignored

 * Disk storage service improvements
   * Allow for storage backend errors - optionally ignore if configured (e.g.: for NFS)
   * Refresh every n seconds
   * Look at improving start-up time using parallel directory exploration (queue w/ goroutines?)

 * Alternative storage services
   * AWS S3 backed storage, with database containing metadata to avoid having to download tracks from S3 every start-up

 * Optionally allow storage backend instances to be named (and use this in error/log messages).

 * Keep player at top of page when scrolling long list
   * Less important now that the filename -> album regexp code is working (previously a directory of 1000s of MP3s would show up as one long album) 

 * Better favicon for showing in browser title bar

 * Structured logging using the standard log package or something from https://blog.logrocket.com/5-structured-logging-packages-for-go/ ?

 * Listen to media hotkeys in browser (and optionally allow this to be disabled?)
 * Stop playing music when the computer is being suspended/hibernated

 * Bug: Still can't skip on long tracks when using remote minimediaserver (not local one) (PROBABLY FIXED by Chrome+range changes)
   * This seems to be a Firefox HTML audio player thing, rather than a golang HTTP server thing.
   * Would it help to be able to support fetching data ranges on a track? Instead of streaming the response with echo .Stream()?
   * JS code could pre-fetch, to populate the browser's cache?
   * Completely override fetching the data in JS, and just give the player a buffer of data instead of a URL?

 * Tags:
   * FLAC (Vorbis comments) (DONE, needs unit test coverage)
   * Ogg (Vorbis comments) (DONE, needs unit test coverage)
   * ID3 (DONE, needs unit test coverage)
   * ID3 v2 (DONE, needs unit test coverage)
   * M4A files from iTunes
   * Use for artist + title instead of filenames in playlists (DONE)
   * Generate playlists based off tags (if present) rather than file location - playlist per album (DONE)

 * Include byte ranges in the HTTP logs
   * Might help debug weird pauses seen in Chrome on macOS, where my RPi takes 3+ minutes to send most of the data???
   * Possibly related issue: https://issues.chromium.org/issues/40942481
   * I think it was a Chrome issue - I haven't seen this with newer Chromes on macOS.
   * Although maybe something is off? My RPi 4 seems to wedge periodically with network buffer errors and I can't ssh into it. Is minimediaserver causing it?

 * Try out on mobile phone. Need media queries to adjust layout for smaller screens or readability?

 * Memory profiling
 * Coverage
 * Function documentation

 * Track data JSON blob - fetch that via API rather than including in generated HTML data
   * with OpenAPI schema and validation in golang code
 * Optimize track storage in media server (*Track instead of Track?)
