using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using UnityEngine;

public class TCPClient
{
  public readonly string host;
  public readonly int port;

  private TcpClient tcpClient;
  private Thread listenThread;

  // Our message handler, which almost certainly should push messages onto a thread-safe queue
  // so it can be handled on the main game thread.
  public delegate void MessageHandler(byte[] bytes);
  private MessageHandler handler;

  public TCPClient(string host, int port, MessageHandler handler)
  {
    this.host = host;
    this.port = port;
    this.handler = handler;
  }

  public void ConnectAndListen()
  {
    tcpClient = new TcpClient(host, port);

    listenThread = new Thread(new ThreadStart(ListenLoop));
    listenThread.IsBackground = true;
    listenThread.Start();
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
          // TODO: Optimize this to not allocate for each message received, and instead only expand the buffer and store
          //  the number of bytes written, resetting this effective "buffer tail" after each message is complete.
          byte[] bytes = new byte[msgLength];
          bytesRead = 0;
          while (bytesRead < msgLength)
          {
            bytesRead += stream.Read(bytes, bytesRead, msgLength - bytesRead);
          }
          handler(bytes);
        }
      }
    }
    catch (SocketException se)
    {
      Debug.Log("Error: " + se);
    }
  }

  // Send assumes that the rawMsg passed is already length-prefixed. eg, using FlatBuffers's `FinishSizePrefixed...` functions.
  public void Send(byte[] rawMsg)
  {
    NetworkStream stream = tcpClient.GetStream();
    if (stream.CanWrite)
    {
      stream.Write(rawMsg, 0, rawMsg.Length);
    }
  }

  public void Close()
  {
    if (listenThread != null)
    {
      listenThread.Abort();
    }
    if (tcpClient != null)
    {
      tcpClient.Close();
    }
  }
}
