import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiClient {
  final String baseUrl;
  ApiClient({required this.baseUrl});

  Future<Map<String, dynamic>> evaluate(List<String> hole, List<String> community) async {
    final res = await http.post(Uri.parse('$baseUrl/evaluate'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'hole': hole, 'community': community}));
    return jsonDecode(res.body) as Map<String, dynamic>;
  }

  Future<Map<String, dynamic>> compare(Map<String, dynamic> hand1, Map<String, dynamic> hand2) async {
    final res = await http.post(Uri.parse('$baseUrl/compare'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'hand1': hand1, 'hand2': hand2}));
    return jsonDecode(res.body) as Map<String, dynamic>;
  }

  Future<Map<String, dynamic>> simulate(Map<String, dynamic> req) async {
    final res = await http.post(Uri.parse('$baseUrl/simulate'),
        headers: {'Content-Type': 'application/json'}, body: jsonEncode(req));
    return jsonDecode(res.body) as Map<String, dynamic>;
  }
}
