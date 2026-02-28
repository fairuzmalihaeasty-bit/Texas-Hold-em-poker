import 'package:flutter/material.dart';
import '../api/client.dart';

class SimulateScreen extends StatefulWidget {
  @override
  State<SimulateScreen> createState() => _SimulateScreenState();
}

class _SimulateScreenState extends State<SimulateScreen> {
  final _heroHole = TextEditingController(text: 'As,Kd');
  final _community = TextEditingController(text: '2d,7c,Jd');
  final _numPlayers = TextEditingController(text: '2');
  final _iterations = TextEditingController(text: '1000');
  String _output = '';
  final api = ApiClient(baseUrl: 'http://10.0.2.2:8080');

  void _submit() async {
    final hero = _heroHole.text.split(',').map((s)=>s.trim()).where((s)=>s.isNotEmpty).toList();
    final community = _community.text.split(',').map((s)=>s.trim()).where((s)=>s.isNotEmpty).toList();
    final numPlayers = int.tryParse(_numPlayers.text) ?? 2;
    final iterations = int.tryParse(_iterations.text) ?? 1000;
    setState(() => _output = 'Running...');
    try {
      final req = {
        'hero': {'hole': hero},
        'community_known': community,
        'num_players': numPlayers,
        'iterations': iterations,
        'concurrency': 4,
      };
      final res = await api.simulate(req);
      setState(() => _output = res.toString());
    } catch (e) { setState(() => _output = 'Error: $e'); }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.all(12),
      child: Column(children: [
        Text('Hero hole (2 cards)'), TextField(controller: _heroHole), SizedBox(height:6),
        Text('Known community (0-5)'), TextField(controller: _community), SizedBox(height:6),
        Row(children: [Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children:[Text('Players'), TextField(controller:_numPlayers)])), SizedBox(width:8), Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children:[Text('Iterations'), TextField(controller:_iterations)]))]),
        SizedBox(height:12), ElevatedButton(onPressed: _submit, child: Text('Simulate')),
        SizedBox(height:12), Expanded(child: SingleChildScrollView(child: Text(_output))),
      ]),
    );
  }
}
