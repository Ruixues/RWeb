import 'package:web_socket_channel/io.dart';

class RWebClient {
  final String serverAddress;
  IOWebSocketChannel conn;
  Map<String, dynamic> bindFunction;
  RWebClient({this.serverAddress}) {}
  void onData(dynamic data) {
    // 开始解析
  }
  void onError() {
    Future.delayed(Duration(seconds: 3), connect);
  }

  Future<bool> connect() async {
    conn = IOWebSocketChannel.connect("ws://localhost:1234");
    conn.stream.listen(onData, onError: onError);
    return true;
  }
}
