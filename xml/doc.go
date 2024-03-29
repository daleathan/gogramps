/*
Package xml provides structs for representing a Gramps XML file and helpers.

Information about the Gramps XML format can be found here:
http://www.gramps-project.org/wiki/index.php?title=GRAMPS_XML#Parsing_Gramps_XML_file

The DTD file is the best documentation I've found so far for what the various
fields mean: http://gramps-project.org/xml/1.5.0/grampsxml.dtd. Note that we
only parse a single version of the XML: 1.5.0 (the most recently released).

The fields and structs in this file are, for the most part, named after fields
in the XML. Notable exceptions:
        * The Database struct represents the XML file as a whole.
        * CamelCase is used instead of lowercase (e.g. ChildOf vs. childof)
        * Repeated elements have plural field names (e.g. Names vs name)
        * Elements that just group other repeated elements are collapsed (e.g.  the tags element)
*/
package xml
