package checkstyle

import (
	"encoding/xml"
)

// CSResult represents a checkstyle XML result.
//
//	<?xml version="1.0" encoding="utf-8"?>
//	<checkstyle version="4.3">
//		<file name="filename">
//			<error line="1" column="3" severity="error" message="msg" source="src" />
//			<error line="2" column="9" severity="error" message="msg" source="src" />
//		</file>
//	</checkstyle>
type CSResult struct {
	XMLName xml.Name  `xml:"checkstyle"`
	Version string    `xml:"version,attr"`
	Files   []*CSFile `xml:"file,omitempty"`
}

// CSFile represents a checkstyle XML file element.
type CSFile struct {
	Name   string     `xml:"name,attr"`
	Errors []*CSError `xml:"error"`
}

// CSError represents a checkstyle XML error element.
type CSError struct {
	Line     int    `xml:"line,attr"`
	Column   int    `xml:"column,attr"`
	Message  string `xml:"message,attr"`
	Severity string `xml:"severity,attr,omitempty"`
	Source   string `xml:"source,attr,omitempty"`
}
