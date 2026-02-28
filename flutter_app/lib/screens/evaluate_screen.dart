import 'package:flutter/material.dart';
import '../api/client.dart';

class EvaluateScreen extends StatefulWidget {
  @override
  State<EvaluateScreen> createState() => _EvaluateScreenState();
}

class _EvaluateScreenState extends State<EvaluateScreen> {
  final _holeController = TextEditingController(text: 'As,Kd');
  final _communityController = TextEditingController(text: '2d,7c,Jd');
  String _output = '';
  final api = ApiClient(baseUrl: 'http://10.0.2.2:8080');

  void _submit() async {
    final hole = _holeController.text.split(',').map((s) => s.trim()).where((s)=>s.isNotEmpty).toList();
    final community = _communityController.text.split(',').map((s) => s.trim()).where((s)=>s.isNotEmpty).toList();
    setState(() => _output = 'Loading...');
    try {
      final res = await api.evaluate(hole, community);
      setState(() => _output = res.toString());
    } catch (e) {
      setState(() => _output = 'Error: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.all(12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text('Hole cards (comma-separated)'),
          TextField(controller: _holeController),
          SizedBox(height: 8),
          Text('Community cards (0-5, comma-separated)'),
          TextField(controller: _communityController),
          SizedBox(height: 12),
          ElevatedButton(onPressed: _submit, child: Text('Evaluate')),
          SizedBox(height: 12),
          Expanded(child: SingleChildScrollView(child: Text(_output))),
        ],
      ),
    );
  }
}
