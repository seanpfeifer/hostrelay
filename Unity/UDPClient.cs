using System;
using System.Collections;
using System.Collections.Generic;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using UnityEngine;

public class UDPClient
{
  public readonly string host;
  public readonly int port;

  private IPEndPoint udpServerEndpoint;
  private UdpClient udpClient;
  private Thread listenThread;

  public UDPClient(string host, int port)
  {
    this.host = host;
    this.port = port;
    this.udpServerEndpoint = new IPEndPoint(IPAddress.Any, port);
  }

  public void ConnectAndListen()
  {
    udpClient = new UdpClient(host, port);

    listenThread = new Thread(new ThreadStart(ListenLoop));
    listenThread.IsBackground = true;
    listenThread.Start();

    SendMessage("UDP client here!");
  }

  private void ListenLoop()
  {
    try
    {
      while (true)
      {
        byte[] bytes = udpClient.Receive(ref udpServerEndpoint);
        Debug.Log("Received UDP: " + Encoding.UTF8.GetString(bytes));
      }
    }
    catch (SocketException se)
    {
      Debug.Log("Error: " + se);
    }
  }

  public void SendMessage(string msg)
  {
    if (udpClient == null)
    {
      return;
    }

    byte[] msgBytes = Encoding.UTF8.GetBytes(msg);
    udpClient.Send(msgBytes, msgBytes.Length);
  }

  public void Close()
  {
    listenThread.Interrupt();
    udpClient.Close();
  }
}
