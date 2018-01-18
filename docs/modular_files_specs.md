# Modular Files Specs (Work in Progress)

**Example Tree**

* [`kitties`]
    * [`group_0`] (kitty group 0)
        * `config.json` (config file for kitty group 0)
        * [`body`]
            * `0_outline.png` (image outline file)
            * `0_area.png` (optional - image area file)
            * `0_config.json` (option config file for image)
            * `...`
        * [`head`]
            * `...`
    * [`group_1`] (kitty group 1)
        * `...`
* [`skins`] (different cat skins, can have patterns. For: head, ears, body, tail)
    * [`0.png`]
    * [`1.png`]
* [`colors`] (colors/gradients for eyes/nose)
    
## Details

### Kitty part format

This is now split into 3 files.

* `/kitties/group_0/head/0_outline.png`

    This is essentially a black outline of the shape of `head_0`, and any parts of the head where we don't want to change the color of.
    
    Anything other than the outline should be transparent. Parts of the `0_outline.png` file can have varying alpha, as the outline will go over the `0_area.png` 
    
    
* `/kitties/group_0/head/0_area.png`