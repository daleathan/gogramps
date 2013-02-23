/*
Identity parses a gramps XML file and then re-serializes it.

This can be used as an example program for writing a program that
modifies the parsed XML before re-serializing it.

It can also be used to verify that the parsing is working correctly.  Given a
gramps export, running identity on the file, loading it into Gramps and
re-exporting it should generate an identical file to the first export.
Unfortunately, the load into Gramps is necessary because Gramps canonicalizes
the ordering of elements and attributes.

Example:
identity -in=foo.gramps -out=bar.gramps
*/
package documentation
