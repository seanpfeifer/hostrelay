using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class GameClient : MonoBehaviour
{
  public string host = "localhost";
  public int tcpPort = 8080;
  public int udpPort = 8585;
  private TCPClient tcpClient;
  private UDPClient udpClient;
  // Start is called before the first frame update
  void Start()
  {
    tcpClient = new TCPClient(host, tcpPort);
    udpClient = new UDPClient(host, udpPort);

    tcpClient.ConnectAndListen();
    udpClient.ConnectAndListen();
  }
}
