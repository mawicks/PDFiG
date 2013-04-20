Portable Document Format In Go (PDFiG)
======================================

Generate and manipulate PDF files using the Go programming language

This project has the goal of being a full-featured PDF library with an
emphasis on both reading *and* writing PDF files.  Currently, only
writing is supported, but the architecture has been designed with
reading in mind.  For example, object references read from one PDF
file will be transparently translated and renumbered for use in
another PDF file.  The project has the goal of being able to read
arbitrary objects from one PDF file and write them to another PDF
file.  It has an additional goal of being able to make append-only,
incremental revisions to existing PDF files.  The design of the
Portable Document Format specifically provides for such revisions.

The API currently has two major classes: `pdf.File` and
`pdf.Document`.  The `pdf.File` class deals with low-level file
structure.  It is meant for low-level manipulation of PDF files and
will not necessarily produce a file that can be read by a PDF reader
without imposing additional document structure.  The `pdf.Document`
class is a high-level API that uses a `pdf.File` object to generate
usable PDF files.  PDF files may be constructed using only a
`pdf.Document` object, without regard to the underlying `pdf.File`
object.

Although the library is very much a work in progress, it is currently
quite usable for producing PDF files.
