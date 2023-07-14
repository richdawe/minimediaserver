# Design for Track Naming and Playlists

## Original Design

This was pretty simple, based on the directory structure. Any `.m3u` playlist files were ignored. The playlist was built on the directories and filenames alone.

E.g.: for a base directory of `home/rdawe/Music/cds` and an album like:

```bash
$ find /home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose | sort -n
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/01.Stacked_Actors.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/02.BreakOut.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/03.Learn_to_Fly.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/04.Gimme_Stitches.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/05.Generator.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/06.Aurora.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/07.Live-In_Skin.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/08.Next_Year.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/09.Headwires.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/10.Aint_It_The_Life.flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/11.M.I.A..flac
/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/There_Is_Nothing_Left_To_Lose.m3u
```

The artist would be `Foo_Fighters`, the album would be `There_Is_Nothing_Left_To_Lose`, and the tracks would be `01.Stacked_Actors`, etc. Given this, a playlist called `Foo_Fighters :: There_Is_Nothing_Left_To_Lose` would be constructed containing tracks `01.Stacked_Actors`, etc.

This scheme works pretty well if you've structured your music library like this. This includes e.g.: an iTunes music library.

This scheme is also pretty simple to implement. There's no ambiguity about which playlist a track belongs to. If it's in a directory, it is in the playlist for that directory.

This scheme does not work so well for older music libraries like MP3s that may have been extracted with a flat structure. It also ignores any Vorbis or ID3 tags that the files may have. In many cases, it also does not work very for albums that have a different artist per track (e.g.: compilations or mixes).

Here are some examples from my MP3 library:

```bash
$ find /home/rdawe/Music/mp3/ -name 'good*' | sort -n | head -n 5
/home/rdawe/Music/mp3/good looking - logical progression (disc 01) (01) - ltj bukem - demon's theme.mp3
/home/rdawe/Music/mp3/good looking - logical progression (disc 01) (02) - chameleon - links.mp3
/home/rdawe/Music/mp3/good looking - logical progression (disc 01) (03) - ltj bukem - music.mp3
/home/rdawe/Music/mp3/good looking - logical progression (disc 01) (04) - pfm - one & only.mp3
/home/rdawe/Music/mp3/good looking - logical progression (disc 01) (05) - aquarius & tayla - bringing me down.mp3

$ find /home/rdawe/Music/mp3/ -name '*monkey kong*' | sort -n | head -n 5
/home/rdawe/Music/mp3/a - a vs. monkey kong (01) - for starters.mp3
/home/rdawe/Music/mp3/a - a vs. monkey kong (02) - monkey kong.mp3
/home/rdawe/Music/mp3/a - a vs. monkey kong (03) - a.mp3
/home/rdawe/Music/mp3/a - a vs. monkey kong (04) - old folks.mp3
/home/rdawe/Music/mp3/a - a vs. monkey kong (05) - hopper jonnus fang.mp3
```

## New Design

### Goals

 * Use Vorbis or ID3 tags (if present) to determine artist, album, etc.

 * Generate playlists based on tracks rather than file location, so that a directory full of music is split into multiple playlists.

 * Support determining artist, album, etc. from the filename using regular expressions, to match e.g.: my MP3 library.

 * If it's not possible to determine the artist, album, etc. using one method, fall back to the next one:
    * 1st: Tags
    * 2nd: Regular expression matches on filename.
    * 3rd: Directory + filename
    * Note: In some situations we may need to use tags for some fields, but directory when e.g.: album artist is missing (see sections below).

### Virtual Playlist Locations

Each playlist has a location, which is a unique name for it. For playlists based on file location, this is the directory name, including the base directory for the storage instance. E.g.: `/home/rdawe/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose` for the album "There Is Nothing Left To Lose" by the Foo Fighters.

What should we use for playlists based on tags or regex matches? Let's use fake URLs, e.g.: `tags:/basepath/artist/album` and `regex:/basepath/artist/album`. A more concrete example might look like `regex:/home/rdawe/Music/mp3/artist/album`.

Note that there is a playlist for each album by an artist. There is no support for a playlist covering a selection of songs, or multiple albums by an artist. Supporting that is not currently a goal.

The playlist ID is generated using the location. The playlist ID should be stable - i.e.: the same across restarts of minimediaserver.

### Albums with Multiple Artists

Albums may have multiple artists. (I'm using the term "artist" here to cover author, composer, performer.) E.g.: an orchestra playing pieces by multiple composers. Or a compilation by multiple artists. In this case the album may have an overall "album artist" and then per-track artists. It is possible that the artist for the entire album is also the artist for one of the tracks.

For ID3 tags, there is a way to represent multiple artists, using TPE1, TPE2 and TPE3 comments. For Vorbis comments for Ogg and FLAC, there does not seem to be a standard tag. And in any case, the relevant tags may not be present - some heuristics may be required.

If it's possible to determine the "album artist" from the tags on a track, then the tracks can be grouped using (album artist, album).

For FLAC files with no album artist tags, the directory structure might indicate that it's a multi-author one (e.g.: "Various Artists/Album/track1.flac, etc.). In that case, we would want to use the directory name as the album artist, but could use the tags to determine the artist for the track. A reasonable heuristic here is that if directory name != artist from tags, then it's a multi-artist album. Artist names should be compared case-insensitively.

The files may have a CDDB tag, containing the unique ID for the album. This may be used as a reliable way of grouping tracks.

### Grouping Tracks into a Playlist

Note: .m3u files are ignored by this design. Perhaps a future version of this design will consider using them. They are ignored in preference of using the audio files' tags as the source of truth.

To build a playlist, the media server needs to:

 * Find audio files ("tracks")
 * Examine each track to determine its artist, album, title, and album artist ("annotate" the track information)
 * Group tracks by (album artist, album) - one playlist per (album artist, album)

When tags are present, we should prefer grouping using the most unique fields available, in this order of preference:

 * `tags:/basepath/albumId/album` (e.g.: using CDDB database ID for the album)
 * `tags:/basepath/albumArtist/album` (e.g.: when it's a multi-artist album)
 * `tags:/basepath/artist/album`

TODO: Issues:

 * FLAC Vorbis tags seem to be missing overall artist tag (e.g.: "Various") - need some heuristics to figure out actual artist? Can use the CDDB tag to match albums too.
 * Probably need post-processing step for playlist to handle multi-artist album
 * Mermaid-format diagram of processing pipeline for DiskStorageService, since it's not straightforward anymore ;)
