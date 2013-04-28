<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
    <xsl:output method="xml" indent="yes"/>
    <xsl:strip-space elements="*"/>
    <!-- identity transform -->
    <xsl:template match="@*|node()">
        <xsl:copy>
            <xsl:apply-templates select="@*|node()"/>
        </xsl:copy>
    </xsl:template>
    <!-- sort node children by their `id` attributes
    <xsl:template match="node()">
        <xsl:copy>
            <xsl:apply-templates select="@*"/>
            <xsl:for-each select="node()">
                <xsl:sort select="@id" order="ascending"/>
                <xsl:apply-templates select="."/>
            </xsl:for-each>
        </xsl:copy>
    </xsl:template> -->
</xsl:stylesheet>
