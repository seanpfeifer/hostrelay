using System;
using System.Collections;
using System.Collections.Generic;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using UnityEngine;

public class TCPClient
{
  public readonly string host;
  public readonly int port;
  private TcpClient tcpClient;
  private Thread listenThread;

  public TCPClient(string host, int port)
  {
    this.host = host;
    this.port = port;
  }

  public void ConnectAndListen()
  {
    tcpClient = new TcpClient(host, port);

    listenThread = new Thread(new ThreadStart(ListenLoop));
    listenThread.IsBackground = true;
    listenThread.Start();

    SendMessage("TCP client here!");
  }

  private void ListenLoop()
  {
    try
    {
      Byte[] prefixBytes = new Byte[4];
      int msgLength;
      int bytesRead;

      using (NetworkStream stream = tcpClient.GetStream())
      {
        while (true)
        {
          // Read the prefix. This will tell us the length of the incoming message.
          bytesRead = 0;
          while (bytesRead < 4)
          {
            bytesRead += stream.Read(prefixBytes, bytesRead, 4 - bytesRead);
          }
          msgLength = BitConverter.ToInt32(prefixBytes, 0);

          // Now read the message
          // TODO: Optimize this to not allocate for each message received
          Byte[] bytes = new Byte[msgLength];
          bytesRead = 0;
          while (bytesRead < msgLength)
          {
            bytesRead += stream.Read(bytes, bytesRead, msgLength - bytesRead);
          }
          Debug.Log("Received: " + Encoding.UTF8.GetString(bytes));
        }
      }
    }
    catch (SocketException se)
    {
      Debug.Log("Error: " + se);
    }
  }

  public void SendMessage(string msg)
  {
    if (tcpClient == null)
    {
      return;
    }

    NetworkStream stream = tcpClient.GetStream();
    if (stream.CanWrite)
    {
      byte[] msgBytes = Encoding.UTF8.GetBytes(msg);
      byte[] prefix = BitConverter.GetBytes(msgBytes.Length);

      stream.Write(prefix, 0, prefix.Length);
      stream.Write(msgBytes, 0, msgBytes.Length);
    }
  }

  public void Close()
  {
    listenThread.Interrupt();
    tcpClient.Close();
  }
}
