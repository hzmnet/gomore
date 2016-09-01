# Gomore

This library is inspired by the Less library and css extension. Less is not more in this case. I took one look at the documentation for Less and thought: Jesus Christ. 

Gomore is much less than Less, it is aimed to be a fast and simple css parser that adds a few useful functions like variables, and mixins, but doesn't try to add everything *and* the kitchen sink.

It is also written in go-lang, not JS (node.js), which makes it much more friendly and fast to use, as it's possible to crosscompile executables.

Gomore is intended to be a library that I will use with my other project, Hyde: https://github.com/jasonknight/hyde

## Status

In development, pre-alpha.

## Examples

```css
$big = 14px;
.some-class {
    font-size: $big;
    color: blue;
    margin-bottom: 1em;
}
.some-other-class {
    .some-class(font-size,color);
    margin-bottom: 2em;
}

```

All selectors are also functions that return their attributes. If you do not specify which ones you want, all will be returned. This is how mixing in is accomplished. Obviously `$big` is a variable. 

There is no special way to define a snippet, just give it a name that won't be used.

```css
.some-class-that-is-never-applied {
    border: 1px solid black;
    border-radius: 5px;
    padding: 5px;
}

.a-used-class {
    .some-class-that-is-never-applied();
    border-radius: 3px;
}

// Should produce:

.a-used-class {
    border: 1px solid black;
    padding: 5px;
    border-radius: 3px;
}

```

The design of the system is such that duplication is in most cases avoided. Notice that `border-radius` is not duplicated.

