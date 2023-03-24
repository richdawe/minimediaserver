# minimediaserver

A mini media server, which provides a web interface for browsing and listening to your music.

`example-config.json` is a basic example configuration. Copy it to `$HOME/..minimediaserver.json` for it to be used by the media server.

E.g.: with a config like this:

```json
{
	"host": "*",
	"port":	"1337",

	"storageServices": [
		{
			"type": "nullStorage"
		},
		{
			"type": "diskStorage",
			"path": "$HOME/Music/cds"
		}
	]
}
```

It can read a music library with a layout like this:

```bash
$ find ~/Music/cds | grep -i foo_fighters
/home/username/Music/cds/Foo_Fighters
/home/username/Music/cds/Foo_Fighters/Foo_Fighters
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/01.This_Is_a_Call.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/12.Exhausted.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/09.For_All_The_Cows.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/08.Oh,_George.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/11.Wattershed.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/02.Ill_Stick_Around.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/Foo_Fighters.m3u
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/06.Floaty.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/05.Good_Grief.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/04.Alone+Easy_Target.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/03.Big_Me.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/07.Weenie_Beenie.flac
/home/username/Music/cds/Foo_Fighters/Foo_Fighters/10.X-Static.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/08.Next_Year.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/01.Stacked_Actors.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/05.Generator.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/11.M.I.A..flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/There_Is_Nothing_Left_To_Lose.m3u
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/10.Aint_It_The_Life.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/02.BreakOut.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/04.Gimme_Stitches.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/03.Learn_to_Fly.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/07.Live-In_Skin.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/09.Headwires.flac
/home/username/Music/cds/Foo_Fighters/There_Is_Nothing_Left_To_Lose/06.Aurora.flac
...
```