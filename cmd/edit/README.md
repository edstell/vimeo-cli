# edit
This shell script wraps the `vimeo Videos Edit` method simplifying editing the settings of a video.

A template config file has been provided (called `config.template`). Full documentation of fields available for editing can be found in the [API Documentation](https://developer.vimeo.com/api/reference/videos#edit_video).

## Usage
Before running, make sure you have the `vimeo` executable in the same folder as edit.sh, and you have set your access token (by editing `edit.sh`).

To edit a video:
1. Create a configuration file using the template provided. 
2. Alter the fields you want to update and remove any fields you don't wish to alter.
3. Execute the script like so:

```
./edit.sh <video_id> <config_path>
```
