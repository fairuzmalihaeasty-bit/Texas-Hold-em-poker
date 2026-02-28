import 'package:flutter/material.dart';
import 'screens/evaluate_screen.dart';
import 'screens/compare_screen.dart';
import 'screens/simulate_screen.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatefulWidget {
  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  int _index = 0;
  final screens = [EvaluateScreen(), CompareScreen(), SimulateScreen()];

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Texas Hold\'em Client',
      theme: ThemeData(primarySwatch: Colors.indigo),
      home: Scaffold(
        appBar: AppBar(title: Text('Texas Hold\'em')),
        body: screens[_index],
        bottomNavigationBar: BottomNavigationBar(
          currentIndex: _index,
          onTap: (i) => setState(() => _index = i),
          items: [
            BottomNavigationBarItem(icon: Icon(Icons.card_giftcard), label: 'Evaluate'),
            BottomNavigationBarItem(icon: Icon(Icons.compare), label: 'Compare'),
            BottomNavigationBarItem(icon: Icon(Icons.show_chart), label: 'Simulate'),
          ],
        ),
      ),
    );
  }
}
