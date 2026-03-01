import 'package:flutter/material.dart';
import '../api/client.dart';
import '../config.dart';

class CompareScreen extends StatefulWidget {
  @override
  State<CompareScreen> createState() => _CompareScreenState();
}

class _CompareScreenState extends State<CompareScreen> {
  final _h1Hole = TextEditingController(text: 'As,Ah');
  final _h1Community = TextEditingController(text: 'Kd,Qd,Jd,2s,3c');
  final _h2Hole = TextEditingController(text: 'Ks,Kh');
  final _h2Community = TextEditingController(text: 'Kd,Qd,Jd,2s,3c');
  String _output = '';
  final api = ApiClient(baseUrl: BASE_URL);

  void _submit() async {
    final hand1 = {
      'hole': _h1Hole.text.split(',').map((s)=>s.trim()).where((s)=>s.isNotEmpty).toList(),
      'community': _h1Community.text.split(',').map((s)=>s.trim()).where((s)=>s.isNotEmpty).toList(),
    };
    final hand2 = {
      'hole': _h2Hole.text.split(',').map((s)=>s.trim()).where((s)=>s.isNotEmpty).toList(),
      'community': _h2Community.text.split(',').map((s)=>s.trim()).where((s)=>s.isNotEmpty).toList(),
    };
    setState(() => _output = 'Loading...');
    try {
      final res = await api.compare(hand1, hand2);
      setState(() => _output = res.toString());
    } catch (e) { setState(() => _output = 'Error: $e'); }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.all(12),
      child: Column(children: [
        Text('Hand 1 - Hole'), TextField(controller: _h1Hole), SizedBox(height: 6),
        Text('Hand 1 - Community'), TextField(controller: _h1Community), SizedBox(height: 8),
        Text('Hand 2 - Hole'), TextField(controller: _h2Hole), SizedBox(height: 6),
        Text('Hand 2 - Community'), TextField(controller: _h2Community), SizedBox(height: 12),
        ElevatedButton(onPressed: _submit, child: Text('Compare')),
        SizedBox(height: 12), Expanded(child: SingleChildScrollView(child: Text(_output))),
      ],),
    );
  }
}
