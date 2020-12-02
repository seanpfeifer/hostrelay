using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using UnityEngine;

public class UDPClient
{
  public readonly string host;
  public readonly int port;

  private IPEndPoint udpServerEndpoint;
  private UdpClient udpClient;
  private Thread listenThread;

  // Our message handler, which almost certainly should push messages onto a thread-safe queue
  // so it can be handled on the main game thread.
  public delegate void MessageHandler(byte[] bytes);
  private MessageHandler handler;

  public UDPClient(string host, int port, MessageHandler handler)
  {
    this.host = host;
    this.port = port;
    this.udpServerEndpoint = new IPEndPoint(IPAddress.Any, port);
    this.handler = handler;
  }

  public void ConnectAndListen()
  {
    udpClient = new UdpClient(host, port);

    listenThread = new Thread(new ThreadStart(ListenLoop));
    listenThread.IsBackground = true;
    listenThread.Start();
  }

  private void ListenLoop()
  {
    try
    {
      while (true)
      {
        byte[] bytes = udpClient.Receive(ref udpServerEndpoint);
        handler(bytes);
      }
    }
    catch (SocketException se)
    {
      Debug.Log("Error: " + se);
    }
  }

  public void Send(byte[] rawMsg)
  {
    udpClient.Send(rawMsg, rawMsg.Length);
  }

  public void Close()
  {
    if (listenThread != null)
    {
      listenThread.Abort();
    }
    if (udpClient != null)
    {
      udpClient.Close();
    }
  }
}
